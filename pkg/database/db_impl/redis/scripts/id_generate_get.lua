local key = KEYS[1];
local keepAliveKey = KEYS[2];
local maxID = tonumber(ARGV[1]);
local timeout = tonumber(ARGV[2]);

-- init bitmap
local e = redis.call("EXISTS", key);
if e == 0 then
    redis.call("SETBIT", key, maxID, 0);
else
    local len = redis.call("STRLEN", key);
    if len * 8 < maxID then
        redis.call("SETBIT", key, maxID, 0);
    end
end

-- now
local t = redis.call('TIME')
local now = tonumber(t[1])

-- free timeout id
local timeoutIds = redis.call("ZRANGEBYSCORE", keepAliveKey, 0, now - timeout);
for _, id in ipairs(timeoutIds) do
    redis.call("SETBIT", key, id, 0);
    redis.call("ZREM", keepAliveKey, id);
end

-- get id
local id = redis.call("BITPOS", key, 0, 0, maxID, "BIT");
if id == -1 then
    return redis.error_reply("NOT_ENOUGH_ID");
end

-- mark id as used
redis.call("SETBIT", key, id, 1);
redis.call("ZADD", keepAliveKey, now, id);

return id;