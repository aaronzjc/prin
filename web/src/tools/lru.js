const CACHE_KEY = "qrcode_lru"
const MAX_SIZE = 10

export function get() {
    let data = localStorage.getItem(CACHE_KEY)
    if (!data) {
        return []
    }
    let logs = JSON.parse(data)
    return logs
}

export function add(item) {
    let data = localStorage.getItem(CACHE_KEY)
    let logs = []
    if (data) {
        logs = JSON.parse(data)
    }
    item = item.trim()
    let res = []
    res.push(item)
    for (let i = 0; i < logs.length; i++) {
        if (logs[i] != item) {
            res.push(logs[i])
        }
    }
    res = res.slice(0, MAX_SIZE-1)
    localStorage.setItem(CACHE_KEY, JSON.stringify(res))
    return res
}

export function remove(idx) {
    let data = localStorage.getItem(CACHE_KEY)
    let logs = []
    if (data) {
        logs = JSON.parse(data)
    }

    logs.splice(idx, 1)
    localStorage.setItem(CACHE_KEY, JSON.stringify(logs))
}