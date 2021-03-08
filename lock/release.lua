local key=KEYS[1]
local newOwner=ARGV[1]
if redis.call('exists',key)==0 then
    return 0
end
if redis.call('hget', key,'owner') == newOwner then
	    local new_ref_count= redis.call('hincrby',key,'ref_count',-1)
	    --无人引用，释放锁
	    if new_ref_count<=0 then
	        return redis.call('del',key)
	    end
	    return new_ref_count
	else
	    return 0
end