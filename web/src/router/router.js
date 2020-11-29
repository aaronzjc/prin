import { createRouter, createWebHashHistory } from 'vue-router'

const Qrcode = () => import("../components/Qrcode")
const Time = () => import("../components/Time")
const Coder = () => import("../components/Coder")
const Cert = () => import("../components/Cert")
const Iptable = () => import("../components/Iptable")

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
    {
        path: "/iptable",
        name: "iptable",
        title: "IPTable",
        component: Iptable
    }
];
const router = createRouter({
    history: createWebHashHistory(),
    routes: baseRoutes.concat(mainRoutes)
})

export {mainRoutes}

export default router