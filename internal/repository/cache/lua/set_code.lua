-- 验证码在redis上的key
local key = KEYS[1]
-- 验证次数，一个验证码最多重复三次，cntKey记录的是还可以验证几次
local cntKey = key..":cnt"
-- 验证码
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在，但是没有过期时间
    return -2
    -- -2是key不存在，ttl < 540 是发了一个验证码，已经超过一分钟了，可以重新发送
elseif ttl == -2  or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 已经发送了一个验证码，但是还不到一分钟
    return -1
end