/// <reference types="vite/client" />
/// <reference types="vue-i18n" />

interface ImportMetaEnv {
    readonly VITE_APPSTORE_BUILD?: string
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}

declare module '*.vue' {
    import type {DefineComponent} from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
}
