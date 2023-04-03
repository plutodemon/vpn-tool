<template>

  <el-menu
      :collapse=$store.state.collapse
      router
      class="elm"
  >
    <el-menu-item >
      <el-icon @click="$store.state.collapse=!$store.state.collapse">
        <Expand/>
      </el-icon>
      <span class="aside_top">VPN管理</span>
    </el-menu-item>

    <el-menu-item index="users">
      <el-icon>
        <User/>
      </el-icon>
      <span class="item">账户</span>
    </el-menu-item>

    <el-menu-item index="adduser">
      <el-icon>
        <CirclePlus/>
      </el-icon>
      <span style="margin-left: 7px">添加账户</span>
    </el-menu-item>


    <el-menu-item index="nets">
      <el-icon>
        <Link/>
      </el-icon>
      <span class="item">网络</span>
    </el-menu-item>

    <el-menu-item index="addnet">
      <el-icon>
        <Monitor/>
      </el-icon>
      <span style="margin-left: 7px">添加网络</span>
    </el-menu-item>

    <el-menu-item index="getspeed">
      <el-icon>
        <Odometer/>
      </el-icon>
      <span style="margin-left: 7px">系统监控</span>
    </el-menu-item>
  </el-menu>
</template>

<script lang="ts">
import {defineComponent, inject} from "vue";
import {
  Monitor,
  View,
  Location,
  Setting,
  Expand,
  Odometer,
  Link,
  User,
  CirclePlus
} from '@element-plus/icons-vue'

export default defineComponent({
  name: "Aside",
  components: {
    Monitor,
    View,
    Location,
    Setting,
    Odometer,
    Link,
    User,
    CirclePlus
  },
  data() {
    return {
      //collapse: false,
    }
  },
  methods: {}
})
</script>

<style scoped>
.aside_top {
  font-size: 20px;
}

.elm {
  background: rgba(255, 255, 255, 0);
  min-height: calc(97vh);
  display: flex;
  flex-flow: column nowrap;
}

.item {
  margin-left: 20px;
  font-size: 15px;
}
</style>