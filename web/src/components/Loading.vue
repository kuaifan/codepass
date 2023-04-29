<template>
    <svg viewBox="25 25 50 50" class="loading">
        <circle cx="50" cy="50" r="20" fill="none" stroke-width="5" stroke-miterlimit="10" class="circle"></circle>
    </svg>
</template>

<style lang="less" scoped>
.loading {
    transform-origin: center center;
    width: 30px;
    height: 30px;
    max-width: 100%;
    max-height: 100%;
    margin: auto;
    overflow: hidden;
    animation: loading-rotate 2s linear infinite;

    .circle {
        fill: none;
        stroke-dasharray: 1, 200;
        stroke-dashoffset: 0;
        stroke-linecap: round;
        animation: path-dash 1.5s ease-in-out infinite, path-color 6s ease-in-out infinite;
    }
}

@keyframes loading-rotate {
    to {
        transform: rotate(1turn)
    }
}

@keyframes path-dash {
    0% {
        stroke-dasharray: 1, 200;
        stroke-dashoffset: 0
    }

    50% {
        stroke-dasharray: 89, 200;
        stroke-dashoffset: -35
    }

    to {
        stroke-dasharray: 89, 200;
        stroke-dashoffset: -124
    }
}

@keyframes path-color {
    0%, to {
        stroke: #d62d20
    }

    40% {
        stroke: #0057e7
    }

    66% {
        stroke: #008744
    }

    80%, 90% {
        stroke: #ffa700
    }
}
</style>
<script>
import {defineComponent} from "vue";

export default defineComponent({
    name: "Loading",
})
</script>
