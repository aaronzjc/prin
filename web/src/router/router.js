import { createRouter, createWebHashHistory } from 'vue-router'

const Qrcode = () => import("../components/pieces/Qrcode")
const Time = () => import("../components/pieces/Time")
const Coder = () => import("../components/pieces/Coder")
const Cert = () => import("../components/pieces/Cert")

const baseRoutes = [
    {
        path: "/",
        name: "default",
        redirect: "qrcode",
    }
];
const mainRoutes = [
    {
        path: "/qrcode",
        name: "qrcode",
        title: "二维码",
        component: Qrcode
    },
    {
        path: "/time",
        name: "time",
        title: "时间",
        component: Time
    },
    {
        path: "/coder",
        name: "coder",
        title: "编解码",
        component: Coder
    },
    {
        path: "/cert",
        name: "cert",
        title: "本地证书",
        component: Cert
    },
];
const router = createRouter({
    history: createWebHashHistory(),
    routes: baseRoutes.concat(mainRoutes)
})

export {mainRoutes}

export default router