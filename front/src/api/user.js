import request from "@/utils/request";

export function getAllEnterprise(query) {
    return request({
        url: '/ui/createUser',
        method: 'get',
        params: query
    })
}