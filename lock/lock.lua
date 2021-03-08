local key=KEYS[1]
local newOwner=ARGV[1]
local repeatable=ARGV[2]
local expireSecond=ARGV[3]
if redis.call('exists',key)==0 then
    redis.call('hset',key,'owner',newOwner,'repeat',repeatable,'ref_count',1)
    return 1
end
if redis.call('hget', key,'owner') == newOwner then
        if redis.call('hget',key,'repeat')==0 then
            return 0
        end
	    local new_ref_count= redis.call('hincrby',key,'ref_count',1)
	    --刷新过期时间（保证就算网络有延迟，从锁设定后的过期时间与代码保持一致）
	    redis.call('expire',key,expireSecond)
	    return new_ref_count
	else
	    return 0
end