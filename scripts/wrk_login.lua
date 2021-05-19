math.randomseed(os.time())

function file_exists(file)
  local f = io.open(file, "rb")
  if f then f:close() end
  return f ~= nil
end

function lines_from(file)
  if not file_exists(file) then return {} end
  lines = {}
  for line in io.lines(file) do
    lines[#lines + 1] = line
  end
  return lines
end

local seed_file = './scripts/seed.txt'
seed_data_lines = lines_from(seed_file)

request = function()
    index = math.random(#seed_data_lines)
    line = seed_data_lines[index]

    wrk.method = "POST"
    wrk.body   = string.format("username=%s&password=%s", line, line)
    wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
    return wrk.format(nil, "/user/login")
end
