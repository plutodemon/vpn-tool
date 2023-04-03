<template>
  <div class="box">
    <div class="form">
      <el-form :model="ipAddress" status-icon ref="ruleForm" label-position="top">

        <el-form-item label="网络号" prop="ip" class="form-item" :rules="rules.ip">
          <el-input v-model="ipAddress.ip" class="form-input"></el-input>
        </el-form-item>


        <el-form-item label="掩码" prop="netMask" class="form-item" :rules="rules.ip">
          <el-input v-model="ipAddress.netMask" class="form-input"></el-input>
        </el-form-item>

        <el-form-item label="网关" prop="netGateway" class="form-item" :rules="rules.ip">
          <el-input v-model="ipAddress.netGateway" class="form-input"></el-input>
        </el-form-item>

        <el-form-item label="DNS地址" prop="dnsAddress" class="form-item" :rules="rules.ip">
          <el-input v-model="ipAddress.dnsAddress" class="form-input"></el-input>
        </el-form-item>

        <el-form-item label="最大设备" prop="maxAllow" class="form-item">
          <el-input v-model.number="ipAddress.maxAllow" class="form-input"></el-input>
        </el-form-item>


        <el-form-item>
          <el-button @click="submitForm" class="form-button">提交</el-button>
          <el-button @click="resetForm('ruleForm')" class="form-button">重置</el-button>
        </el-form-item>

      </el-form>
    </div>
  </div>

</template>

<script>
import {defineComponent, markRaw} from "vue";
import ipInput from "@/components/ip/ipInput";
import {ElMessage, ElMessageBox} from "element-plus";
import {Link} from "@element-plus/icons-vue";
import axios from "axios";

export default defineComponent({
  name: "AddNet",
  components: {
    ipInput
  },
  data() {
    return {
      ipAddress: {
        ip: '',
        netMask: '',
        netGateway: '',
        dnsAddress: '',
        maxAllow: undefined
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
    };
  },
  methods: {
    resetForm(formName) {
      this.$refs[formName].resetFields();
    },
    submitForm() {

      //进行表单验证
      this.$refs.ruleForm.validate(valid => {
        //表单验证失败
        if (!valid) {
          //提示语
          this.$message("请正确填写");
          return false;
        }
        ElMessageBox.confirm(
            '即将创建新网络，确认继续？',
            '提示',
            {
              confirmButtonText: '确认',
              cancelButtonText: '返回',
              type: 'success',
              icon: markRaw(Link),
            }
        )
            .then(() => {
              axios.post("/api/ui/createNet", this.ipAddress).then(res => {
                console.log(res)
                if (res.status === 200) {
                  ElMessage({
                    type: 'success',
                    message: '创建成功',
                  })
                } else {
                  ElMessage({
                    type: 'danger',
                    message: '创建失败',
                  })
                }
              })
            })
            .catch(() => {
              ElMessage({
                type: 'info',
                message: '未创建',
              })
            })
      })

    }
  }
})
</script>

<style scoped>
.box {
  display: flex;
  flex-direction: row;
  justify-content: space-around;
}

.form {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 700px;
  height: 600px;
  margin-top: 20px;
  margin-left: 35px;
  margin-right: 35px;

  border-radius: 10px;
  background: rgba(238, 233, 233, 0.25);
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

.form-item {
  width: 400px;
  margin-top: 15px;
  margin-bottom: 15px;
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

.form-input {
  width: 380px;
  background: rgba(238, 233, 233, 0.25);
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
}
</style>