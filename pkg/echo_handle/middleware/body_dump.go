package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"silo/pkg/errs"
	"silo/pkg/types/isotime"
	"silo/pkg/utils"

	"io"
	"net"
	"net/http"
	"reflect"
	"time"

	"log/slog"

	"dario.cat/mergo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestLog .
type RequestLog struct {
	UserID uint64 `json:"user_id,omitempty"`
	// CreatedAt is time recorded before next middleware/handler is executed.
	CreatedAt isotime.ISOTime `json:"created_at,omitempty"`
	// Latency is duration it took to execute rest of the handler chain (next(c) call).
	Latency float64 `json:"latency,omitempty"`
	// Protocol is request protocol (i.e. `HTTP/1.1` or `HTTP/2`)
	Protocol string `json:"protocol,omitempty"`
	// RemoteIP is request remote IP. See `echo.Context.RealIP()` for implementation details.
	RemoteIP string `json:"remote_ip,omitempty"`
	// Host is request host value (i.e. `example.com`)
	// Host string
	// Method is request method value (i.e. `GET` etc)
	Method string `json:"method,omitempty"`
	// URI is request URI (i.e. `/list?lang=en&page=1`)
	URI string `json:"uri,omitempty"`
	// URIPath is request URI path part (i.e. `/list`)
	URIPath string `json:"uri_path,omitempty"`
	// RoutePath is route path part to which request was matched to (i.e. `/user/:id`)
	RoutePath string `json:"route_path,omitempty"`
	// RequestID is request ID from request `X-Request-ID` header or response if request did not have value.
	RequestID string `json:"request_id,omitempty"`
	// Referer is request referer values.
	Referer string `json:"referer,omitempty"`
	// UserAgent is request user agent values.
	UserAgent string `json:"user_agent,omitempty"`
	// Status is response status code. Then handler returns an echo.HTTPError then code from there.
	Status int `json:"status,omitempty"`
	// Error is error returned from executed handler chain.
	Error         string `json:"error,omitempty"`
	Stack         string `json:"stack,omitempty"`
	err           error  `json:"-"`
	isCustomError bool   `json:"-"`
	// ContentLength is content length header value. Note: this value could be different from actual request body size
	// as it could be spoofed etc.
	// ContentLength string
	RequestBody any `json:"request,omitempty"`
	// ResponseBody is response body value. Note: it can be too big.
	ResponseBody any `json:"response,omitempty"`
	// Headers are list of headers from request. Note: request can contain more than one header with same value so slice
	// of values is been logger for each given header.
	// Note: header values are converted to canonical form with http.CanonicalHeaderKey as this how request parser converts header
	// names to. For example, the canonical key for "accept-encoding" is "Accept-Encoding".
	Headers map[string][]string `json:"headers,omitempty"`
	// QueryParams are list of query parameters from request URI. Note: request can contain more than one query parameter
	// with same name so slice of values is been logger for each given query param name.
	QueryParams map[string][]string `json:"query_params,omitempty"`
	// FormValues are list of form values from request body+URI. Note: request can contain more than one form value with
	// same name so slice of values is been logger for each given form value name.
	FormValues map[string][]string `json:"form_values,omitempty"`
}

// RequestLoggerConfig is configuration for Request Logger middleware.
type RequestLoggerConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// BeforeNextFunc defines a function that is called before next middleware or handler is called in chain.
	BeforeNextFunc func(c echo.Context)
	// LogValuesFunc defines a function that is called with values extracted by logger from request/response.
	// Mandatory.
	LogValuesFunc func(c echo.Context, v *RequestLog) error
	// LogLatency instructs logger to record duration it took to execute rest of the handler chain (next(c) call).
	LogLatency bool
	// LogProtocol instructs logger to extract request protocol (i.e. `HTTP/1.1` or `HTTP/2`)
	LogProtocol bool
	// LogRemoteIP instructs logger to extract request remote IP. See `echo.Context.RealIP()` for implementation details.
	LogRemoteIP bool
	// LogHost instructs logger to extract request host value (i.e. `example.com`)
	// LogHost bool
	// LogMethod instructs logger to extract request method value (i.e. `GET` etc)
	LogMethod bool
	// LogURI instructs logger to extract request URI (i.e. `/list?lang=en&page=1`)
	LogURI bool
	// LogURIPath instructs logger to extract request URI path part (i.e. `/list`)
	LogURIPath bool
	// LogRoutePath instructs logger to extract route path part to which request was matched to (i.e. `/user/:id`)
	LogRoutePath bool
	// LogRequestID instructs logger to extract request ID from request `X-Request-ID` header or response if request did not have value.
	LogRequestID bool
	// LogReferer instructs logger to extract request referer values.
	LogReferer bool
	// LogUserAgent instructs logger to extract request user agent values.
	LogUserAgent bool
	// LogStatus instructs logger to extract response status code. If handler chain returns an echo.HTTPError,
	// the status code is extracted from the echo.HTTPError returned
	LogStatus bool
	// LogError instructs logger to extract error returned from executed handler chain.
	LogError       bool
	LogRequestBody bool
	// LogResponseBody instructs logger to extract response body. Note: it can be to big.
	LogResponseBody bool
	// LogHeaders instructs logger to extract given list of headers from request. Note: request can contain more than
	// one header with same value so slice of values is been logger for each given header.
	//
	// Note: header values are converted to canonical form with http.CanonicalHeaderKey as this how request parser converts header
	// names to. For example, the canonical key for "accept-encoding" is "Accept-Encoding".
	LogHeaders []string
	// LogQueryParams instructs logger to extract given list of query parameters from request URI. Note: request can
	// contain more than one query parameter with same name so slice of values is been logger for each given query param name.
	LogQueryParams []string
	// LogFormValues instructs logger to extract given list of form values from request body+URI. Note: request can
	// contain more than one form value with same name so slice of values is been logger for each given form value name.
	LogFormValues []string
}

// NewBodyDump logger
func NewBodyDump() echo.MiddlewareFunc {
	return RequestLoggerWithConfig(RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, values *RequestLog) error {
			// handle panic
			defer func() {
				if r := recover(); r != nil {
					slog.Error("[api] panic", "panic", r)
				}
			}()

			// check logskipper
			s := c.Get(CtxKeyLogSkipper)
			skipper, ok := s.(*LogSkipperConfig)
			skipFullLog := ok && skipper.ShouldSkipFullLog(c)
			if !skipFullLog || values.Status >= http.StatusBadRequest {
				skipRequestLog := ok && skipper.ShouldSkipRequestLog(c)
				if skipRequestLog {
					values.RequestBody = "*"
				}

				args := []any{}
				v := reflect.ValueOf(values).Elem()
				t := reflect.TypeOf(values).Elem()
				for i := 0; i < v.NumField(); i++ {
					if v.Field(i).IsValid() && !v.Field(i).IsZero() && v.Field(i).CanInterface() {
						tag := utils.StructFieldJSONTag(t.Field(i))
						if tag != "" {
							args = append(args, tag, v.Field(i).Interface())
						}
					}
				}

				if values.Status >= http.StatusBadRequest {
					// args = append(args, "headers", c.Request().Header)
					slog.Warn("[api]", args...)
				} else if values.Status >= http.StatusInternalServerError {
					// args = append(args, "headers", c.Request().Header)
					slog.Error("[api]", args...)
				} else if values.err != nil {
					if values.isCustomError {
						slog.Warn("[api]", args...)
					} else {
						slog.Error("[api]", args...)
					}
				} else {
					slog.Info("[api]", args...)
				}

			}

			return nil
		},
	})
}

// RequestLoggerWithConfig returns a RequestLogger middleware with config.
func RequestLoggerWithConfig(config RequestLoggerConfig) echo.MiddlewareFunc {
	def := RequestLoggerConfig{
		LogLatency:  true,
		LogProtocol: false,
		LogRemoteIP: true,
		// LogHost:         true,
		LogMethod:       true,
		LogURI:          true,
		LogURIPath:      false,
		LogRoutePath:    false,
		LogRequestID:    false,
		LogReferer:      false,
		LogUserAgent:    true,
		LogStatus:       true,
		LogError:        true,
		LogRequestBody:  true,
		LogResponseBody: true,
	}

	err := mergo.Merge(&def, config)
	if err != nil {
		panic(err)
	}

	mw, err := def.ToMiddleware()
	if err != nil {
		panic(err)
	}
	return mw
}

// ToMiddleware converts RequestLoggerConfig into middleware or returns an error for invalid configuration.
func (config RequestLoggerConfig) ToMiddleware() (echo.MiddlewareFunc, error) {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	if config.LogValuesFunc == nil {
		return nil, errors.New("missing LogValuesFunc callback function for request logger middleware")
	}

	logHeaders := len(config.LogHeaders) > 0
	headers := append([]string(nil), config.LogHeaders...)
	for i, v := range headers {
		headers[i] = http.CanonicalHeaderKey(v)
	}

	logQueryParams := len(config.LogQueryParams) > 0
	logFormValues := len(config.LogFormValues) > 0

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}
			req := c.Request()
			res := c.Response()
			start := isotime.Now()
			v := &RequestLog{
				CreatedAt: start,
			}

			// handle panic
			defer func() {
				if r := recover(); r != nil {
					re, ok := r.(error)
					if ok {
						v.err = re
						v.Error = re.Error()
						v.Stack = errs.GetStackTrace(6, 5)
					} else {
						ce := errs.NewInternalServerError("panic")
						v.err = ce
						v.Error = ce.Error()
						v.Stack = errs.GetStackTrace(6, 5)
					}
					config.LogValuesFunc(c, v)

					c.Error(err)

					return
				}
			}()

			if config.LogProtocol {
				v.Protocol = req.Proto
			}
			if config.LogRemoteIP {
				v.RemoteIP = c.RealIP()
			}
			// if config.LogHost {
			// 	v.Host = req.Host
			// }
			if config.LogMethod {
				v.Method = req.Method
			}
			if config.LogURI {
				v.URI = req.RequestURI
			}
			if config.LogURIPath {
				p := req.URL.Path
				if p == "" {
					p = "/"
				}
				v.URIPath = p
			}
			if config.LogRoutePath {
				v.RoutePath = c.Path()
			}
			if config.LogRequestID {
				id := req.Header.Get(echo.HeaderXRequestID)
				if id == "" {
					id = res.Header().Get(echo.HeaderXRequestID)
				}
				v.RequestID = id
			}
			if config.LogReferer {
				v.Referer = req.Referer()
			}
			if config.LogUserAgent {
				v.UserAgent = req.UserAgent()
			}

			if logHeaders {
				v.Headers = map[string][]string{}
				for _, header := range headers {
					if values, ok := req.Header[header]; ok {
						v.Headers[header] = values
					}
				}
			}
			if logQueryParams {
				queryParams := c.QueryParams()
				v.QueryParams = map[string][]string{}
				for _, param := range config.LogQueryParams {
					if values, ok := queryParams[param]; ok {
						v.QueryParams[param] = values
					}
				}
			}
			if logFormValues {
				v.FormValues = map[string][]string{}
				for _, formValue := range config.LogFormValues {
					if values, ok := req.Form[formValue]; ok {
						v.FormValues[formValue] = values
					}
				}
			}

			// Request
			reqBody := []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = io.ReadAll(c.Request().Body)
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			if config.BeforeNextFunc != nil {
				config.BeforeNextFunc(c)
			}

			// Response
			resBody := new(bytes.Buffer)

			//match routes
			routes := c.Echo().Routes()
			var routeMatched bool
			for _, r := range routes {
				if r.Path == c.Path() && r.Method == req.Method {
					routeMatched = true
					break
				}
			}
			if !routeMatched {
				err = errs.NewNotFoundError("route not found").WithHTTPCode(http.StatusNotFound)
			} else {
				if config.LogResponseBody {
					mw := io.MultiWriter(res.Writer, resBody)
					res.Writer = &bodyDumpResponseWriterLogger{Writer: mw, ResponseWriter: res.Writer}
				}

				err = next(c)
			}

			end := time.Now()

			m := c.Get(CtxKeyMasker)
			masker, hasMasker := m.(*Masker)

			// request
			if hasMasker {
				v.RequestBody = masker.MaskerRequest(string(reqBody))
			} else {
				v.RequestBody = string(reqBody)
			}

			var rb any
			if err := json.Unmarshal(reqBody, &rb); err == nil {
				v.RequestBody = rb
			}

			// response
			if config.LogResponseBody {
				if hasMasker {
					v.ResponseBody = masker.MaskerResponse(resBody.String())
				} else {
					v.ResponseBody = resBody.String()
				}

				var body any
				if err := json.Unmarshal(resBody.Bytes(), &body); err == nil {
					v.ResponseBody = body
				}
			}

			// apply user claims
			// if session, err := common.SessionFromCtx(c); err == nil {
			// 	v.UserID = session.UID
			// }

			if config.LogLatency {
				v.Latency = end.Sub(start.ToTime()).Seconds()
			}

			if config.LogStatus {
				v.Status = res.Status
				if err != nil {
					if httpErr, ok := err.(*echo.HTTPError); ok {
						v.Status = httpErr.Code
					}
				}
			}
			if config.LogError && err != nil {
				v.err = err
				v.Error = err.Error()

				if ce, ok := err.(*errs.CustomError); ok {
					v.isCustomError = true
					v.Stack = ce.StackTrace()
				} else {
					v.isCustomError = false
				}
			}

			config.LogValuesFunc(c, v)

			if err != nil {
				if ce, ok := err.(*errs.CustomError); ok {
					ce.WithoutInternalError()
				}
				c.Error(err)
			}

			return
		}
	}, nil
}

type bodyDumpResponseWriterLogger struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriterLogger) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriterLogger) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriterLogger) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriterLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
