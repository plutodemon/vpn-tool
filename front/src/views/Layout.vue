<template>
  <div class="common-layout">
    <el-container>

      <el-aside width="auto">
        <Aside/>
      </el-aside>

      <el-container>

        <el-header>
          <Head/>
        </el-header>

        <el-main>
          <router-view/>
        </el-main>

      </el-container>
    </el-container>
  </div>
</template>


<script lang="ts">
//侧边栏
import {defineComponent} from 'vue'
import Head from "@/components/body/Head.vue";
import Aside from "@/components/body/Aside.vue";
import Main from "@/components/body/Main.vue";

export default defineComponent({
  name: "Layout",
  components: {
    Head,
    Aside,
    Main
  },
})
</script>


<style>

</style>
