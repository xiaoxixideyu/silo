local key = KEYS[1];
local keepAliveKey = KEYS[2];
local id = tonumber(ARGV[1]);

redis.call("SETBIT", key, id, 0);
redis.call("ZREM", keepAliveKey, id);

return id;