math.randomseed(os.time())
math.random(); math.random(); math.random()

request = function()
    wrk.method = "POST"
    count = math.random(1, 200)

    loginRes = math.random(0, 2)
    if math.ceil(loginRes) == 2 then
        password = "world"
    else
        password = "worldd"
    end

    wrk.body   = "username=hello" .. math.ceil(count) .. "&password=" .. password
    wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
    return wrk.format(nil, "/user/login")
end
