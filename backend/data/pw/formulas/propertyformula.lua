
local percent = 0.009999999
function get_HP(p1,p4,p5,p6)
        if p6==13 then p1=p1-25 end
	return p1 * (100 + p4) * percent + p5 + 0.5
end


function get_MP(p1,p4,p6)
        if p6==13 then p1=p1+25 end
	return p1 * (100 + p4) * percent + 0.5
end


local function enh1(weapon_type, host_prof, str, agi, vita)
	local T = 0
	
	if weapon_type == 1 then
		if host_prof == 4 then
			T = math.floor((str + vita) * (100 / 150) + 0.5)
		else
			T = math.floor(str * (100 / 150) + 0.5)
		end
	elseif weapon_type == 0 or weapon_type == 2 or weapon_type == 3 then
		T = math.floor(agi * (100 / 150) + 0.5)
	end
	
	return T
end

function get_damage(p1, p4, weapon_type, host_prof, str, agi, vita)
	local enh = enh1(weapon_type, host_prof, str, agi, vita)
	
	return math.floor((p1 + 1) * (100 + p4 + enh) * percent + 0.5)
end


function get_damage_magic(p1,p4, p5)
	return (p1 + 1) * (100 + p4 + p5) * percent + 0.5
end


function get_attack(p1,p4)
	return p1 * (100 + p4) * percent + 0.5
end


function get_armor(p1,p3)
	return p1 * (100 + p3) * percent + 0.5
end


function get_attack_degree(p1,p2,p3)

	if p2 < -100 then p2 = -100 end
	return (100 + p2) * 0.01 * p1 + p3
end


function get_defend_degree(p1,p2,p3)
	if p2 < -100 then p2 = -100 end
	return (100 + p2) * 0.01 * p1 + p3
end


function get_penetration(p1,p2,p3)
	local rst = 0
	if p2 < -100 then p2 = -100 end
	
	rst = p1 * ((100 + p2) * percent + 0.00001)
	
	if rst < 0 then 
		rst = 0
	elseif p1 < 0 and rst > 0 then
		rst = 0
	end
	
	rst = rst + p3
	
	return rst
end


function get_resilience_level(p1,p2,p3)
	local rst = 0
	
	if p2 < -100 then p2 = -100 end
	
	rst = p1 * ((100 + p2) * percent + 0.00001)
	
	if rst < 0 then
		rst = 0 
	elseif p1 < 0 and rst > 0 then
		rst = 0
	end
	
	rst = rst + p3
	
	return rst
end
	

function get_crit_rate(p1,agi)
	local p2 = 0
	p2 = agi * 0.05 + 0.001
	p2 = p2 + p1

	if p2 > 100 then p2 = 100 end
	return p2
end


function get_crit_damage(p1)
	return p1 + 200
end


function get_crit_damage_defense (p1)
	
	return p1
end


function get_defense(p1,p3,vit,str)
	local enh = 0
	enh = math.floor((vit * 2 + str * 3)*(100 / 2500)+0.5)
	local point_enh = 0
	point_enh = (vit + str)/4
	
	local T = 0
	T = math.floor((p1 + 1) * (100 + p3 + enh) * percent + 0.5) + point_enh
	if T < 0 then T = 0 end
	return T
end


local function get_resilience(p1,p3,p4,vit,eng)
	local res_enh = 0
	res_enh = math.floor((vit * 2 + eng * 3)*(100 / 2500) + 0.5)
	
	local res_point_enh = 0
	res_point_enh = (vit + eng) / 4
	
	local T = 0
	T = math.floor(p1 * (100 + p3 + res_enh) * percent + 0.5) + res_point_enh + p4
	if T < 0 then T = 0 end
	
	return T
end


function get_defense_magic(p1,p2,p3,p4,p5,p6,p7,p8,p9,p10,p11,vit,eng)
	local T1 = get_resilience(p1,p2,p11,vit,eng)
	local T2 = get_resilience(p3,p4,p11,vit,eng)
	local T3 = get_resilience(p5,p6,p11,vit,eng)
	local T4 = get_resilience(p7,p8,p11,vit,eng)
	local T5 = get_resilience(p9,p10,p11,vit,eng)
	return (T1 + T2 + T3 + T4 +T5) / 5
end


function get_attack_speed(weapon_speed,p1,p2)
	weapon_speed = (weapon_speed - p1) * 20
	
	local attack_speed
	
	attack_speed = weapon_speed * percent * (100 - p2)

	
	if attack_speed < 4 then
		attack_speed = 4
	elseif attack_speed > 300 then
		attack_speed = 300
	end
	
	attack_speed = 20 / attack_speed
	return attack_speed
end


function get_prayspeed(p1)
	if p1 >= 100 then
		p1 = 90
	end
	return p1
end


function get_speed(p1,p2,p3)
	local en_speed = p1
	if en_speed > 3 then
		if en_speed > 5 then
			en_speed = 0;
		else
			en_speed = 3
		end
	end
	
	local sp = 0
	sp = p2 * percent * (p3 + 100) + en_speed
	
	if sp > 15 then
		sp = 15
	elseif sp < 1e-3 then
		sp = 0.1
	end
	
	return sp
end


function get_damage_reduce(p1)
	if p1 > 75 then
		p1 = 75
	end
	return p1
end


function get_magic_damage_reduce(p1)
	return p1
end


function get_invisible_degree(p1)
		return 0
end


function get_anti_invisible_degree(p1)
	return p1
end


function get_anti_defense_degree(p1)
	return p1
end


function get_resistance_degree(p1)
	return p1
end


function get_vigour (p1)
	return p1
end


function get_soul_power (p1)
	return p1
end

function get_base_point (p1)
	return p1
end


function get_PVEdamageAS (paramTable,level)
	local T0 = (paramTable[3] + paramTable[4]) / 2 * (1 - percent * paramTable[9] + 0.0001 * paramTable[9] * paramTable[10])
				* (1 + paramTable[17] / 100) * (1 + paramTable[19] / 4000) * (1 + 3 * paramTable[27] / (paramTable[27] + 300))
				* (1 + paramTable[23] / 10000) *  paramTable[12]
	local T1 = (paramTable[5] + paramTable[6]) / 2 * (1 - percent * paramTable[9] + 0.0001 * paramTable[9] * paramTable[10])
				* (1 + paramTable[17] / 100) * (1 + paramTable[19] / 4000) * (1 + 3 * paramTable[27] / (paramTable[27] + 300))
				* (1 + paramTable[24] / 10000) * (2 / (2 - percent * paramTable[13]))
	if	T0 > T1 then
		return T0
	else
		return T1
	end
end


function get_PVPdamageAS (paramTable,level)
	local T0 = (paramTable[3] + paramTable[4]) / 2 * (1 - percent * paramTable[9] + 0.0001 * paramTable[9] * paramTable[10])
				* (1 + paramTable[17] / 100) * (1 + paramTable[19] / 1000) 
				* (1 + paramTable[23] / 10000) * paramTable[12]
	local T1 = (paramTable[5] + paramTable[6]) / 2 * (1 - percent * paramTable[9] + 0.0001 * paramTable[9] * paramTable[10])
				* (1 + paramTable[17] / 100) * (1 + paramTable[19] / 1000)
				* (1 + paramTable[24] / 10000) * (2 / (2 - percent * paramTable[13]))
	if	T0 > T1 then
		return T0
	else
		return T1
	end
end


function get_PVEdefenseAS (paramTable,level)
	local defense = (paramTable[7] / (paramTable[7] + 40 * level - 25))
	local magic_defense = (paramTable[8] / (paramTable[8] + 40 * level -25))
	if defense > 0.95 then
		defense = 0.95
	end
	if magic_defense > 0.95 then
		magic_defense = 0.95
	end
	
	local HP1 = paramTable[1] / (1 - defense) / (1 - paramTable[25]*percent) * (1 + 0.012 * paramTable[18])
				* (1 + paramTable[19] / 4000) * (1 + paramTable[28] / level)
	local HP2 = paramTable[1] / (1 - magic_defense) / (1 - paramTable[26]*percent) * (1 + 0.012 * paramTable[18])
				* (1 + paramTable[19] / 4000) * (1 + paramTable[28] / level)
	local T = (HP1 + HP2) / 2
	
	return T
end


function get_PVPdefenseAS (paramTable,level)
	local defense = (paramTable[7] / (paramTable[7] + 40 * level - 25))
	local magic_defense = (paramTable[8] / (paramTable[8] + 40 * level -25))
	if defense > 0.95 then
		defense = 0.95
	end
	if magic_defense > 0.95 then
		magic_defense = 0.95
	end
	
	local HP1 = paramTable[1] / (1 - defense) / (1 - paramTable[25]*percent) * (1 + 0.012 * paramTable[18])
				* (1 + paramTable[19] / 1000) * (1 - 0.5 + 0.5 * 2) / (1 - 0.5 + 0.5 * (2 - paramTable[11]*percent))
	local HP2 = paramTable[1] / (1 - magic_defense) / (1 - paramTable[26]*percent) * (1 + 0.012 * paramTable[18])
				* (1 + paramTable[19] / 1000)
	local T = (HP1 + HP2) / 2
	
	return T
end

function get_toplevel (p1)	
	return p1
end





