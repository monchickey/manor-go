local key = KEYS[1]
local error_rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local ttl_second = tonumber(ARGV[3])
local items = {}

for i = 4, #ARGV do
    items[i - 3] = ARGV[i]
end

-- local k1, k2 = unpack(items)
-- redis.debug('k1', k1)
-- redis.debug('k2', k2)

local exists = redis.call('EXISTS', key)

if exists == 1 then
    return redis.call('BF.MADD', key, unpack(items))
else
    -- 初始化布隆过滤器
    redis.call('BF.RESERVE', key, error_rate, capacity)
    redis.call('EXPIRE', key, ttl_second)
    return redis.call('BF.MADD', key, unpack(items))
end