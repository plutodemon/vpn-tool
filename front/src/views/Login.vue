<template>
  <div class="login">
    <img
        src="https://ts1.cn.mm.bing.net/th/id/R-C.cfa1ff7fbcb79dd5bb5df921d437b367?rik=GpvgK310HRcHnA&riu=http%3a%2f%2fwww.obzhi.com%2fwp-content%2fuploads%2f2020%2f03%2fearth.jpg&ehk=77eklhcpE%2faDHBOR2s%2flCGbG4Ustmha%2f5XzVBy70tOU%3d&risl=&pid=ImgRaw&r=0"
        style="width:100% ;height:100%;position: absolute; margin-left: -50%;margin-top: -8px" alt=""/>
    <div class="loginPart">
      <h2>用户登录</h2>
      <el-form ref="form" :model="form">

        <div class="inputElement">
          <el-input v-model="form.email" placeholder="请输入邮箱" prefix-icon="User" auto-complete="off"></el-input>
        </div>

        <div class="inputElement">
          <el-input v-model="form.password" placeholder="请输入密码" prefix-icon="Lock" show-password auto-complete="off"></el-input>
        </div>
        <div class="inputElement">
          <el-input v-model="form.code" placeholder="请输入验证码" prefix-icon="Link" auto-complete="off"></el-input>
        </div>

        <el-button type="primary" round style="width: 100px;" @click="login">登录</el-button>

      </el-form>
    </div>
  </div>
</template>

<script lang="ts">
import {
  User,
  Lock,
  Link
} from '@element-plus/icons-vue'
import {defineComponent} from "vue";
import axios from "axios";
import {ElMessage} from "element-plus";

export default defineComponent({
  name: 'Login',

  components: {
    User,
    Lock,
    Link
  },

  data() {
    return {
      checked: true,
      form: {
        email:undefined,
        password:undefined,
        code:undefined
      }
    };
  },

  methods: {
    login() {
      console.log(this.form)
      axios.post("/api/user/login", this.form).then(res => {
        console.log(res)
        if (res.status === 200) {
          ElMessage({
            type: 'success',
            message: '登录成功',
          })
          this.$router.push("/main")
        } else {
          ElMessage({
            type: 'error',
            message: '登录失败',
          })
        }
      })
    }
  }
})
</script>

<style scoped>
.loginPart {
  position: absolute;
  /*定位方式绝对定位absolute*/
  top: 50%;
  left: 50%;
  /*顶和高同时设置50%实现的是同时水平垂直居中效果*/
  transform: translate(-50%, -50%);
  /*实现块元素百分比下居中*/
  width: 450px;
  padding: 50px;
  background: rgba(0, 0, 0, .5);
  /*背景颜色为黑色，透明度为0.8*/
  box-sizing: border-box;
  /*box-sizing设置盒子模型的解析模式为怪异盒模型，将border和padding划归到width范围内*/
  box-shadow: 0 15px 25px rgba(0, 0, 0, .5);
  /*边框阴影  水平阴影0 垂直阴影15px 模糊25px 颜色黑色透明度0.5*/
  border-radius: 15px;
  /*边框圆角，四个角均为15px*/
}

.loginPart h2 {
  margin: 0 0 30px;
  padding: 0;
  color: #fff;
  text-align: center;
  /*文字居中*/
}

.loginPart .inputElement {
  width: 100%;
  padding: 10px 0;
  font-size: 16px;
  color: #fff;
  letter-spacing: 1px;
  /*字符间的间距1px*/
  margin-bottom: 30px;
  border: none;
  border-bottom: 1px solid #fff;
  outline: none;
  /*outline用于绘制元素周围的线outline：none在这里用途是将输入框的边框的线条使其消失*/
  background-color: transparent;
  /*背景颜色为透明*/
}

.login {
  width: 100%;
  height: 100%;
  background: transparent;
}
</style>
