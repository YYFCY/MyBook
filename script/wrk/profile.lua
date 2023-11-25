wrk.method="GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "PostmanRuntime/7.35.0"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDA4OTgwOTcsIlVpZCI6MTMzLCJVc2VyQWdlbnQiOiJQb3N0bWFuUnVudGltZS83LjM1LjAifQ.VSrlFOudVle6sekJ34IpVMemPVU8_JTGIockHY-o17Y"