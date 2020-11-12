import { createRouter, createWebHashHistory } from 'vue-router'

import Qrcode from "../components/pieces/Qrcode"
import Time from "../components/pieces/Time"

const routes = [
    {
        path: "/",
        name: "default",
        redirect: "qrcode",
    },
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
    }
];
const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export {routes}

export default router