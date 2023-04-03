<template>
  <!--  <el-scrollbar max-height="627px" style="display:flex;">-->
  <div class="all_user">
    <template v-for="item in userList" :key="item.id">

      <el-card shadow="hover" class="user">
        <div style=" margin-top:10px;height: 60px">{{ item.userName }}</div>
        <div style="height: 60px">{{ item.userEmail }}</div>
        <div style="height: 60px">
          <el-button icon="EditPen" circle @click="openDia(item)"/>
          <el-button icon="Delete" circle @click="openDelete(item.id)"/>
        </div>
      </el-card>

    </template>
  </div>


  <el-card class="card" shadow="always" v-show="visible">

    <el-form :model="user" status-icon ref="ruleForm" label-position="top" class="form1">
      <el-form-item label="用户名" prop="userName" class="form-item">
        <el-input v-model="user.userName" class="form-input"></el-input>
      </el-form-item>
      <el-form-item label="邮 箱" prop="userEmail" class="form-item">
        <el-input v-model="user.userEmail" class="form-input"></el-input>
      </el-form-item>
      <el-form-item label="网络编号" prop="netId" class="form-item">
        <el-input v-model.number="user.netId" class="form-input"></el-input>
      </el-form-item>
      <el-form-item label="IP地址" prop="netIp" class="form-item" :rules="rules.ip">
        <el-input v-model="user.netIp" class="form-input"></el-input>
      </el-form-item>
      <el-form-item label="用户等级" prop="userLevel" class="form-item">
        <el-input v-model.number="user.userLevel" class="form-input"></el-input>
      </el-form-item>
      <el-form-item label="密码登录" prop="isAllowed" class="form-item">
        <el-switch
            size="large"
            v-model="user.isAllowed"
            class="mt-2"
            style="margin-left: 24px"
            inline-prompt
            active-icon="Check"
            inactive-icon="Close"
        />
      </el-form-item>

      <el-form-item>
        <el-button @click="submitForm('ruleForm')" class="form-button">提交</el-button>
        <el-button @click="visible=false" class="form-button">退出</el-button>
      </el-form-item>
    </el-form>

  </el-card>

  <!--  </el-scrollbar>-->


</template>

<script>
import {EditPen, Delete, Check, Close, User} from "@element-plus/icons-vue";
import {ElMessage, ElMessageBox} from "element-plus";
import {defineComponent, markRaw} from "vue";
import axios from "axios";
import {mapState} from "vuex";

export default defineComponent({
  name: "Users",
  components: {
    EditPen,
    Delete,
    Check,
    Close
  },
  computed: {
    ...mapState({
      collapse: state => state.collapse
    }),
  },
  data() {
    return {
      currentPage: 1,
      userList: [],
      visible: false,
      user: {
        id: undefined,
        userName: undefined,
        userEmail: undefined,
        isAllowed: undefined,
        netId: undefined,
        netIp: undefined,
        userLevel:undefined
      },
      rules: {
        ip: [
          {required: true, message: '请输入地址', trigger: 'blur'},
          {
            validator: function (rule, value, callback) {
              if (/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/.test(value) === false) {
                callback(new Error("请输入正确地址"));
              } else {
                //校验通过
                callback();
              }
            }, trigger: 'blur'
          },
        ]
      }
    }

  },
  methods: {
    openDia(item) {
      this.visible = true
      this.user = item
    },
    submitForm() {
      ElMessageBox.confirm(
          '即将更新用户，确认继续？',
          '提示',
          {
            confirmButtonText: '确认',
            cancelButtonText: '返回',
            type: 'success',
            icon: markRaw(User),
          }
      )
          .then(() => {
            console.log(this.user)
            axios.post("/api/ui/editUser", this.user).then(res => {
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
      ).then(() => {
        axios.get("/api/ui/deleteUser/?id=" + item).then(res => {
          if (res.status === 200) {
            this.getUsers()
            ElMessage({
              type: 'success',
              message: '删除成功',
            })
          }
        })
      }).catch(() => {
        ElMessage({
          type: 'info',
          message: '未删除',
        })
      })
    },
    getUsers() {
      axios.get("/api/ui/findUser").then(res => {
        this.userList = res.data
      })
    }
  },
  created() {
    this.getUsers()
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
  margin-top: 10px;
  margin-bottom: 10px;
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