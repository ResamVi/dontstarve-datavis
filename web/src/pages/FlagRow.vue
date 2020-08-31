<template>
    <div class="row">
        <img :src="require(`../assets/${this.character}.png`)" class="portrait">
        <div class="container" v-for="i in [0, 1, 2, 3, 4]" :key="i">
            <h5>{{ data[i][0] }}</h5>
            <span :class="flag(data[i][1])"></span>
            <h5>{{ data[i][2] }}%</h5>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        character: String,
        data: Array,
    },

    methods: {
        flag(iso) {
            return `flag-icon flag-icon-${iso.toLowerCase()} size`;
        }
    }
}
</script>

<style scoped>
h5 {
    white-space: nowrap;
    font-size: 1vw;
}
.container {
    display: inline;
    text-align: center;
    margin: 0 5px;
}

.row {
    display: grid;
    grid-template-columns: 16% 16% 16% 16% 16% 16%;
}

.portrait {
    width: 7vw;
    height: 7vw;
    justify-self: center;
    align-self: center;
}

.size {
    font-size: 2em;
}


</style>