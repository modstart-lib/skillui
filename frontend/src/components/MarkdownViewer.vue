<template>
    <div class="markdown-viewer" v-html="renderedHtml"></div>
</template>

<script setup lang="ts">
import { marked } from 'marked';
import { computed } from 'vue';

/**
 * MarkdownViewer 只读 Markdown 渲染组件
 *
 * @prop {string} content - Markdown 内容
 */
const props = defineProps({
    content: {
        type: String,
        default: ''
    }
});

// 配置 marked renderer，让所有链接在新窗口打开
const renderer = new marked.Renderer();
const originalLinkRenderer = renderer.link.bind(renderer);

renderer.link = (token: any) => {
    const html = originalLinkRenderer(token);
    return html.replace('<a', '<a target="_blank" rel="noopener noreferrer"');
};

marked.setOptions({
    renderer: renderer
});

const renderedHtml = computed(() => {
    if (!props.content) return '';
    try {
        return marked(props.content);
    } catch (error) {
        console.error('Markdown 渲染错误:', error);
        return props.content;
    }
});
</script>

<style scoped>
.markdown-viewer {
    line-height: 1.5;
    color: #374151;
    font-size: 0.875rem;
    padding: 0;
    margin: 0;
}

.markdown-viewer :deep(h1) {
    font-size: 1.25rem;
    font-weight: 700;
    margin-top: 1rem;
    margin-bottom: 0.5rem;
}

.markdown-viewer :deep(h2) {
    font-size: 1.125rem;
    font-weight: 700;
    margin-top: 0.75rem;
    margin-bottom: 0.5rem;
}

.markdown-viewer :deep(h3) {
    font-size: 1rem;
    font-weight: 700;
    margin-top: 0.5rem;
    margin-bottom: 0.25rem;
}

.markdown-viewer :deep(p) {
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
}

.markdown-viewer :deep(ul),
.markdown-viewer :deep(ol) {
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
    padding-left: 1.5rem;
}

.markdown-viewer :deep(li) {
    margin-top: 0.25rem;
    margin-bottom: 0.25rem;
}

.markdown-viewer :deep(code) {
    background-color: #f3f4f6;
    padding: 0.125rem 0.25rem;
    border-radius: 0.25rem;
    font-size: 0.875rem;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.markdown-viewer :deep(pre) {
    background-color: #f3f4f6;
    padding: 0.75rem;
    border-radius: 0.25rem;
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
    overflow-x: auto;
}

.markdown-viewer :deep(pre code) {
    background-color: transparent;
    padding: 0;
}

.markdown-viewer :deep(blockquote) {
    border-left: 4px solid #d1d5db;
    padding-left: 1rem;
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
    color: #6b7280;
}

.markdown-viewer :deep(a) {
    color: #3b82f6;
    text-decoration: underline;
}

.markdown-viewer :deep(a:hover) {
    color: #1d4ed8;
}

.markdown-viewer :deep(strong) {
    font-weight: 700;
}

.markdown-viewer :deep(em) {
    font-style: italic;
}

.markdown-viewer :deep(table) {
    border-collapse: collapse;
    width: 100%;
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
}

.markdown-viewer :deep(th),
.markdown-viewer :deep(td) {
    border: 1px solid #d1d5db;
    padding: 0.5rem;
    text-align: left;
}

.markdown-viewer :deep(th) {
    background-color: #f3f4f6;
    font-weight: 600;
}
</style>
