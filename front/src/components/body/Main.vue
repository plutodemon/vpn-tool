<template>
  <div class="all_user">

    <div ref="speedCpu" class="user"/>

    <div ref="speedMem" class="user"/>

    <div ref="speedDisk" class="user"/>

    <div ref="numClient" class="user"/>
    <div class="user"><div class="ps"><h1>包速率</h1>{{ pbps }}<h1>pbps</h1></div></div>
    <div  class="user"><div class="ps"><h1>字节速率</h1>{{ kbps }}<h1>kbps</h1></div></div>
  </div>


  <div class="sp">
    <div ref="pbpsSp" style="height: 300px;width: 300px"/>
    <div ref="kbpsSp" style="height: 300px;width: 300px"/>
  </div>

</template>

<script lang="ts">
import {onMounted, ref, onBeforeMount, onUpdated, nextTick, watchEffect, defineComponent} from 'vue'
import * as echarts from 'echarts';
import {Minus, Plus} from '@element-plus/icons-vue'
import {time} from "echarts";

export default defineComponent({

  setup() {
    const clientNum = ref(0)
    const clientNumSp = ref(0)
    const pbps = ref(0)
    const kbps = ref(0)
    const cpuPer = ref(0)
    const memPer = ref(0)
    const diskPer = ref(0)
    const userNum = ref(0)

    const speedCpu = ref(null)
    const speedMem = ref(null)
    const speedDisk = ref(null)
    const numClient = ref(null)
    const speedNumClient = ref(null)
    const pbpsSp = ref(null)
    const kbpsSp = ref(null)

    const pbpsSpData = ref<Array<number>>([]);
    const kbpsSpData = ref<Array<string>>([]);

    const ws = new WebSocket("ws://localhost:8080/ui/getSpeed")
    //连接打开时触发
    ws.onopen = function (evt) {
      console.log("Connection open ...")
      ws.send("getSpeed")
    }
    //接收到消息时触发
    ws.onmessage = function (evt) {
      const data = JSON.parse(evt.data)
      clientNum.value = data.clientNum
      clientNumSp.value = data.clientNumSp
      pbps.value = data.pbps
      kbps.value = data.kbps
      cpuPer.value = Number(data.cpuPer.toFixed(1))
      memPer.value = data.memPer
      diskPer.value = Number(data.diskPer.toFixed(1))
      userNum.value = data.userNum - data.clientNum

    }
    //连接关闭时触发
    ws.onclose = function (evt) {
      console.log("Connection closed.")
    }

    nextTick(() => {
      let speedCpuChart = echarts.init(speedCpu.value!);
      let speedMemChart = echarts.init(speedMem.value!);
      let speedDiskChart = echarts.init(speedDisk.value!);
      let numClientChart = echarts.init(numClient.value!);
      // let speedNumClientChart = echarts.init(speedNumClient.value!);
      let pbpsSpChart = echarts.init(pbpsSp.value!);
      let kbpsSpChart = echarts.init(kbpsSp.value!);

      watchEffect(() => {
        speedCpuChart.setOption({
          series: [
            {
              name: 'Pressure',
              type: 'gauge',
              progress: {
                show: true
              },
              axisTick: {
                show: false
              },
              detail: {
                valueAnimation: true,
                formatter: '{value}'
              },
              data: [
                {
                  value: cpuPer.value,
                  name: 'CPU'
                }
              ]
            }
          ]
        })
        speedMemChart.setOption({
          series: [
            {
              name: 'Pressure',
              type: 'gauge',
              progress: {
                show: true
              },
              axisTick: {
                show: false
              },
              detail: {
                valueAnimation: true,
                formatter: '{value}'
              },
              data: [
                {
                  value: memPer.value,
                  name: '内存'
                }
              ]
            }
          ]
        })
        speedDiskChart.setOption({
          series: [
            {
              name: 'Pressure',
              type: 'gauge',
              progress: {
                show: true
              },
              axisTick: {
                show: false
              },
              detail: {
                valueAnimation: true,
                formatter: '{value}'
              },
              data: [
                {
                  value: diskPer.value,
                  name: '磁盘'
                }
              ]
            }
          ]
        })

        numClientChart.setOption({
          tooltip: {
            trigger: 'item'
          },

          series: [
            {
              type: 'pie',
              radius: ['40%', '70%'],
              center: ['50%', '70%'],
              // adjust the start angle
              startAngle: 180,

              data: [
                {value: clientNum.value, name: '连接用户'},
                {value: userNum.value, name: '未连接用户'},
                {
                  // make an record to fill the bottom 50%
                  value: clientNum.value + userNum.value,
                  itemStyle: {
                    // stop the chart from rendering this piece
                    color: 'none',
                    decal: {
                      symbol: 'none'
                    }
                  },
                  label: {
                    show: false
                  }
                }
              ]
            }
          ]
        })
        // speedNumClientChart.setOption({})
        pbpsSpChart.setOption({})
        kbpsSpChart.setOption({})
      })
    })
    return {
      clientNum,
      clientNumSp,
      pbps,
      kbps,
      cpuPer,
      memPer,
      diskPer,
      speedCpu,
      speedMem,
      speedDisk,
      numClient,
      speedNumClient,
      pbpsSp,
      kbpsSp,
      pbpsSpData,
      kbpsSpData
    }
  },

})
</script>

<style scoped>
.main {
  background: #97c1ee;
}

.sp {
  display: flex;
  flex-flow: row wrap;
  margin-top: 0;
  justify-content: space-around;
  max-height: 88vh;
  position: fixed;
  z-index: 3;
  overflow-y: auto;
}

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
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  width: 400px;
  height: 300px;
  margin-bottom: 20px;
  margin-left: 15px;
  margin-right: 25px;
}

.ps{
  padding-top: 80px;
}
</style>