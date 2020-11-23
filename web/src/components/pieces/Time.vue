<template>
<div class="columns" id="time">
    <div class="column">
        <article class="message is-dark">
            <div class="message-body">
                当前时间 <span class="tag is-info is-light">{{ state.current.time }}</span>
                时间戳 <span class="tag is-info is-light">{{ state.current.timestamp }}</span>
                <span class="tag is-info refresh" @click="fetchCurrentTime">刷新一下</span>
            </div>
        </article>
        <div class="field is-horizontal">
            <div class="field-body">
            
            <div class="field">
            <flat-pickr
                placeholder="输入时间"
                class="input"
                v-model="state.time"
                :config="{dateFormat: 'Y-m-d H:i', enableTime: true}">
            </flat-pickr>
            </div>
            <div class="field is-grouped is-justify-content-center">
            <p class="control">
                <a class="button is-link is-light" @click="toTimestamp">
                转时间戳
                </a>
            </p>
            <p class="control">
                <a class="button is-white">
                = . = 
                </a>
            </p>
            <p class="control">
                <a class="button is-link is-light" @click="toTime">
                转日期
                </a>
            </p>
            </div>
            <div class="field">
            <input class="input" type="text" placeholder="时间戳" v-model="state.timestamp">
            </div>
            </div>
        </div>
    </div>
</div>
</template>

<script>
import { onMounted, reactive } from 'vue'
import flatPickr from 'vue-flatpickr-component';
import 'flatpickr/dist/flatpickr.css';

export default {
    name: "Time",
    setup() {
        const state = reactive({
            current: {
                time: "",
                timestamp: ""
            },
            time: "",
            timestamp: ""
        })

        function timestampToTime(stamp) {
            let dt = new Date(stamp * 1000);
            var year = dt.getFullYear();
            var month = dt.getMonth() + 1;
            var day = dt.getDate();
            var hours = dt.getHours();
            var minutes = dt.getMinutes();
            minutes = minutes < 10 ? '0' + minutes : minutes;
            var seconds = dt.getSeconds();
            seconds = seconds < 10 ? '0' + seconds : seconds;
            return year + "-" + month + "-" + day + " " + hours + ":" + minutes + ":" + seconds;
        }

        function fetchCurrentTime() {
            let dt = new Date();
            state.current.timestamp = Math.floor(dt.getTime() / 1000)
            state.current.time = timestampToTime(dt.getTime() / 1000)
        }

        function toTimestamp() {
            if (state.time == "") {
                state.timestamp = "";
                return
            }
            console.log(state)
            state.timestamp = Math.floor(Date.parse(state.time) / 1000);
        }

        function toTime() {
            if (state.timestamp <= 0) {
                return
            }

            state.time = timestampToTime(state.timestamp)
        }

        onMounted(fetchCurrentTime)

        return {
            state,
            fetchCurrentTime,
            toTimestamp,
            toTime
        }
    },
    components: {
        flatPickr
    }
}
</script>

<style lang="scss">
#time {
    .refresh {
        cursor: pointer;
    }
}
</style>