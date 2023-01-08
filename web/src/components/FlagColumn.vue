<template>
    <div class="container">
        <img :src="imageSrc" class="portrait">
        <div class="column" v-for="(entry, i) in data" :key="i">
            <span class="date">{{ entry["name"] }}</span>
            <span :class="flag(entry['data'][0][0])"></span> 
            <span :class="flag(entry['data'][1][0])"></span> 
            <span :class="flag(entry['data'][2][0])"></span> 
            <span :class="flag(entry['data'][3][0])"></span> 
            <span :class="flag(entry['data'][4][0])"></span> 
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
            return `fi fi-${iso.toLowerCase()} margin`;
        }
    },

    computed: {
        imageSrc() {
            return new URL(`../assets/${this.character}.png`, import.meta.url).href;
        }
    }
}
</script>

<style scoped>
.date {
    font-weight: bold;
    text-align: center
    
}
.container {
    display: grid;
    grid-template-columns: 7vw repeat(12, 60px);
    margin-top: 10px;
    padding-bottom: 5px;
}

.column {
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    align-items: center;
    
}

.margin {
    margin: 2px 0px;
}

.portrait {
    width: 7vw;
    height: 7vw;
    justify-self: center;
    align-self: center;
}
</style>
