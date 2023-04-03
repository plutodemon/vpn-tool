<template>
  <div class="all_user">
    <template v-for="speed in speedList">

      <el-card shadow="hover" class="user">

        <div style=" margin-top:10px;height: 60px">{{ speed.ip }}</div>

        <div style=" margin-top:1px;height: 60px">包速率
          <el-icon>
            <Sort/>
          </el-icon>
          {{ speed.pbps }}
        </div>

        <div style="height: 60px">字节速率
          <el-icon>
            <Sort/>
          </el-icon>
          {{ speed.kbps }}
        </div>
      </el-card>
    </template>
  </div>
</template>

<script>

import {defineComponent, nextTick, watchEffect, ref, reactive, onUnmounted} from "vue";
import axios from "axios";
import {Sort} from "@element-plus/icons-vue";


export default defineComponent({
  name: "Speed",
  components: {
    Sort
  },
  setup() {
    let speedList = ref([])
    nextTick(() => {

      const ws = new WebSocket("ws://localhost:8080/ui/getClientSpeed")
      //连接打开时触发
      ws.onopen = function (evt) {
        console.log("Connection open (getClientSpeed) ...")
        ws.send("getSpeed")
      }
      watchEffect(() => {
        //接收到消息时触发
        ws.onmessage = function (evt) {
          speedList.value = JSON.parse(evt.data)
          console.log("evt.data:", evt.data)
          console.log("this.speedList:", speedList)
        }
      })

      //连接关闭时触发
      ws.onclose = function (evt) {
        console.log("Connection closed (getClientSpeed).")
      }
    })

    /*function getSpeed() {
      axios.get("/api/ui/getClientSpeed").then(res => {
        speedList.value = res.data
        console.log(speedList)
      })
    }

    let t = setInterval(getSpeed, 3000)
    onUnmounted(() => {
      clearInterval(t)
    })*/
    return {
      speedList
    }
  }

})

</script>

<style scoped>
.all_user {
  display: flex;
  flex-flow: row wrap;
  margin-top: 0;
  justify-content: space-around;
  max-height: 88vh;
  position: fixed;
  z-index: 3;
  overflow-y: auto;
}

.user {
  border-radius: 10px;
  background: rgba(238, 233, 233, 0.25);
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  width: 240px;
  height: 200px;
  margin-bottom: 15px;
  margin-left: 15px;
  margin-right: 25px;
}
</style>