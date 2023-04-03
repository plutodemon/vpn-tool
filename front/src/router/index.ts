import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Layout from "@/views/Layout.vue";
// @ts-ignore
import Users from "@/views/user/Users.vue";
import AddUser from "@/views/user/AddUser.vue";
import Nets from "@/views/net/Nets.vue";
import AddNet from "@/views/net/AddNet.vue";
import Main from "@/components/body/Main.vue";
import Login from "@/views/Login.vue";
// @ts-ignore
import Speed from "@/views/controller/Speed"


const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'login',
    component: Login,
  },
  {
    path: '/main',
    name: '主页',
    component: Layout,
    redirect:'/main',
    children:[
      {
        path: '/users',

        name: '账户',
        component: Users,
      },
      {
        path: '/nets',
        name: '网络',
        component: Nets
      },
      {
        path: '/main',
        component: Main
      },
      {
        path: '/adduser',
        name: '添加账户',
        component: AddUser,
      },
      {
        path: '/addnet',
        name: '添加网络',
        component: AddNet,
      },
      {
        path: '/getspeed',
        name: '系统监控',
        component: Speed,
      },
    ]
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
