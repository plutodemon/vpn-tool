import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import store from './store/store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/theme-chalk/dark/css-vars.css'
import VueClipboard from 'vue-clipboard2'

const app = createApp(App)
app.use(store).use(router).use(VueClipboard).use(ElementPlus, {locale: zhCn,})
    .mount('#app')
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

