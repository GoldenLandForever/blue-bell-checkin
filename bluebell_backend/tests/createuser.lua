-- 用户数据结构
local users = {}
local tokens = {}
local completed = 0
local total_users = 100
local base_url = "https://localhost:8081/api/v1"

-- 用户注册请求
function create_user(id)
    -- 生成随机密码
    
    -- 用户信息
    local user = {
        username = "testuser_" .. tostring(id),
        password = "testuser_" .. tostring(id),
    }
    
    return wrk.format("POST", base_url .. "/signup", 
        {
            ["Content-Type"] = "application/json"
        }, 
        wrk.json(user)
    )
end

-- 登录请求（获取token）
function login_user(username, password)
    return wrk.format("POST", base_url .. "/login", 
        {
            ["Content-Type"] = "application/json"
        }, 
        wrk.json({username = username, password = password})
    )
end

-- 初始化函数
function init(args)
    -- 生成所有用户注册请求
    for id = 1, total_users do
        table.insert(users, {
            id = id,
            register_request = create_user(id),
            login_request = nil,
            token = nil
        })
    end
end

-- 响应处理函数
function response(status, headers, body, request)
    if not request or not request.ctx then return end
    
    if request.ctx.stage == "register" then
        if status == 201 then
            local id = request.ctx.id
            local password = users[id].password
            
            -- 生成登录请求
            users[id].login_request = login_user("testuser_" .. tostring(id), password)
            return users[id].login_request, {stage = "login", id = id}
        else
            print("用户注册失败: " .. body)
        end
    
    elseif request.ctx.stage == "login" then
        if status == 200 then
            -- 解析token
            local success, data = pcall(wrk.json.decode, body)
            if success and data.access_token then
                users[request.ctx.id].token = data.access_token
                completed = completed + 1
                
                -- 显示进度
                print("用户 " .. request.ctx.id .. " 创建成功. 进度: " .. completed .. "/" .. total_users)
                
                -- 所有用户完成时保存token
                if completed == total_users then
                    save_tokens()
                end
            else
                print("Token解析失败: " .. body)
            end
        else
            print("用户登录失败: " .. body)
        end
    end
    
    -- 获取下一个注册请求
    for _, user in ipairs(users) do
        if not user.token then
            if not user.login_request then
                return user.register_request, {stage = "register", id = user.id}
            end
        end
    end
    
    -- 所有请求完成
    return nil
end

-- 保存token到文件
function save_tokens()
    local file = io.open("user_tokens.json", "w")
    if not file then return end
    
    local tokens_array = {}
    for _, user in ipairs(users) do
        if user.token then
            table.insert(tokens_array, user.token)
        end
    end
    
    file:write(wrk.json(tokens_array))
    file:close()
    
    print("\n已完成 " .. #tokens_array .. " 用户的创建, tokens保存到 user_tokens.json")
end