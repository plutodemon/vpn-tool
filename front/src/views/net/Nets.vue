<template>
  <el-scrollbar max-height="627px" style="display:flex;">
    <div class="all_user">
      <template v-for="item in netList" :key="item.id">

        <el-card shadow="hover" class="user">
          <div style=" margin-top:10px;height: 60px">{{ item.ip }}</div>
          <div style="height: 60px">{{ item.mask }}</div>
          <div style="height: 60px">
            <el-button icon="EditPen" circle @click="openDia(item)"/>
            <el-button icon="Delete" circle @click="openDelete(item.ID)"/>
          </div>
        </el-card>

      </template>
    </div>

    <el-card class="card" shadow="always" v-show="visible">

      <el-form :model="net" status-icon ref="ruleForm" label-position="top" class="form1">
        <el-form-item label="IP" prop="ip" class="form-item">
          <el-input v-model="net.ip" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="掩码" prop="mask" class="form-item">
          <el-input v-model="net.mask" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="网关" prop="netGateway" class="form-item">
          <el-input v-model="net.netGateway" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="DNS服务器地址" prop="dnsAddress" class="form-item">
          <el-input v-model="net.dnsAddress" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="最大设备数" prop="maxAllow" class="form-item">
          <el-input v-model="net.maxAllow" class="form-input"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button @click="submitForm('ruleForm')" class="form-button">提交</el-button>
          <el-button @click="visible=false" class="form-button">退出</el-button>
        </el-form-item>
      </el-form>

    </el-card>

  </el-scrollbar>


</template>

<script>
import {EditPen, Delete, Check, Close, Link} from "@element-plus/icons-vue";
import {ElMessage, ElMessageBox} from "element-plus";
import {defineComponent, markRaw} from "vue";
import axios from "axios";

export default defineComponent({
  name: "Nets",
  components: {
    EditPen,
    Delete,
    Check,
    Close
  },
  data() {
    return {
      netList: [],
      visible: false,
      net: {
        ID: undefined,
        ip: undefined,
        mask: undefined,
        netGateway: undefined,
        dnsAddress: undefined,
        maxAllow: undefined
      }
    }
  },
  methods: {
    openDia(item) {
      this.visible = true
      this.net = item
    },
    submitForm() {
      ElMessageBox.confirm(
          '即将更新网络，确认继续？',
          '提示',
          {
            confirmButtonText: '确认',
            cancelButtonText: '返回',
            type: 'success',
            icon: markRaw(Link),
          }
      )
          .then(() => {
            console.log(this.net)
            axios.post("/api/ui/editNet", this.net).then(res => {
              console.log(res)
              this.visible = false
              if (res.status === 200) {
                ElMessage({
                  type: 'success',
                  message: '更新成功',
                })
              } else {
                ElMessage({
                  type: 'danger',
                  message: '更新失败',
                })
              }
            })
          })
          .catch(() => {
            ElMessage({
              type: 'info',
              message: '未更新',
            })
          })
    },
    openDelete(item) {
      ElMessageBox.confirm(
          '此用户将被删除，确认继续？',
          '警告',
          {
            confirmButtonText: '确认',
            cancelButtonText: '返回',
            type: 'warning',
            icon: markRaw(Delete),
          }
      )
          .then(() => {
            console.log(this.netList)
            console.log(this.net)
            console.log("item:", item)
            axios.get("/api/ui/deleteNet/?id=" + item).then(res => {
              if (res.status === 200) {
                this.getNets()
                ElMessage({
                  type: 'success',
                  message: '删除成功',
                })
              }
            })
          })
          .catch(() => {
            ElMessage({
              type: 'info',
              message: '未删除',
            })
          })
    },
    getNets() {
      axios.get("/api/ui/findNet").then(res => {
        this.netList = res.data
      })
    }
  },
  created() {
    this.getNets()
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
  margin-bottom: 10px;
  margin-left: 15px;
  margin-right: 25px;
}

.card {
  border-radius: 10px;
  background: rgba(238, 233, 233, 0.5);
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  width: 85%;
  height: 600px;
  display: flex;
  justify-content: center;
  flex-direction: column;
  margin-top: 10px;
  align-items: center;
  align-self: center;
  position: fixed;
  z-index: 5;
}

.form1 {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.form-item {
  width: 400px;
  margin-top: 15px;
  margin-bottom: 15px;
}

.form-input {
  height: 40px;
}

.form-button {
  width: 150px;
  margin-left: 20px;
  margin-right: 20px;
  box-shadow: 0 7px 7px 0 rgba(105, 112, 224, 0.37);
  backdrop-filter: blur(15.5px);
  -webkit-backdrop-filter: blur(15.5px);
  border-radius: 10px;
}
</style>