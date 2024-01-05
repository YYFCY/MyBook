local key = KEYS[1]
local cntKey = key..":cnt"
local expectedCode = ARGV[1]
local cnt = tonumber(redis.call("GET", cntKey))
local code = redis.call("GET", key)

if cnt <= 0 then
    -- 用户一直输错，或者已经用过了
    return -1
end

if code == expectedCode then
    -- 用完不能再用来
    redis.call("set", cntKey, -1)
    return 0
else
    -- 用户输错，可验证次数 -1
    redis.call("decr", cntKey)
    return -2
end