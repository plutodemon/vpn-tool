/* eslint-disable */
import { store } from '@/store/store' // path to store file

declare module '@vue/runtime-core' {
    interface ComponentCustomProperties {
        $store: store;
    }
}
