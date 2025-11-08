local keepAliveKey = KEYS[1];
local id = tonumber(ARGV[1]);

local t = redis.call('TIME')
local now = tonumber(t[1])
-- XX: Only update elements that already exist. Don't add new elements.
redis.call("ZADD", keepAliveKey, 'XX', now, id);

return id;