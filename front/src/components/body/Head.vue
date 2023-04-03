<template>
  <div class="head">
    <div style=" margin-left: 20px">
      <el-breadcrumb separator-icon="ArrowRight">
        <template v-for="(item, index) in breadList">
          <el-breadcrumb-item
              v-if="item.name"
              :key="index"
              :to="item.path"
          >{{ item.name }}
          </el-breadcrumb-item>
        </template>
      </el-breadcrumb>
    </div>
    <div style="margin-right: 20px;">
      <el-switch
          v-model="isDark"
          class="sw_dark"
          size="large"
          inline-prompt
          active-icon="Moon"
          inactive-icon="Sunrise"
      />
      <el-avatar> user</el-avatar>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {useDark, useToggle} from '@vueuse/core'
import {ArrowRight} from '@element-plus/icons-vue'
import {
  Moon,
  Sunrise
} from '@element-plus/icons-vue'
import {mapState} from "vuex";

const isDark = useDark()
const toggleDark = useToggle(isDark)
export default defineComponent({
  name: "Head",
  components: {
    Moon,
    Sunrise,
    ArrowRight
  },
  data() {
    return {
      isDark,
      breadList: [] as any[],
    }
  },
  methods: {
    getMatched() {
      this.breadList = this.$route.matched;
    },
  },
  created() {
    this.getMatched()
  },
  watch: {
    $route(to, from) {
      this.breadList = this.$route.matched;
    }
  },
})
</script>

<style scoped>
.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 50px;
  margin-top: 7px;
  background: rgba(238, 233, 233, 0.25);
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.18);
}

.sw_dark {
  --el-switch-on-color: rgb(154, 136, 136);
  --el-switch-of-color: rgb(178, 119, 119);
  margin-right: 20px;
}
</style>