<template>
    <main> 
        <h1>Welcome!</h1>
        <p>
            DST DataViz records and visualizes current player preferences across regions
            on Klei's <a href="https://store.steampowered.com/app/322330/Dont_Starve_Together/">
            Don't Starve Together</a> game. Data generated 4 minutes ago.
        </p>
        
        <div class="split">
            <div>
                <h1>123456</h1>
                <p>Servers analyzed</p>
            </div>
            <div>
                <h1>12345</h1>
                <p>Players found</p>
            </div>
        </div>

        <h3 class="boxed">
            Count of Characters being played
            <a href="#" class="has-tooltip">[ ? ]
                <span class="tooltip tooltip-top">
                    Any non-DST or non-official characters listed were most likely
                    added to the server via a character mod.
                </span>
            </a>
        </h3>
        <bar-chart :data="characters"></bar-chart>

        <h3 class="boxed">Amount of Servers by Country</h3>
        <bar-chart :data="[['X-Small', 55], ['Small', 27], ['Hard', 22], ['Medium', 7]]"></bar-chart>

        <h3 class="boxed">
            Amount of Players by Country
            <a href="#" class="has-tooltip">[ ? ]
                <span class="tooltip tooltip-top">
                    Players inherit the server's country.
                    For example a Swede, a German and a Dane on a French server count as 3 French players.
                    Klei's API does not give more detail to infer a player's origin.
                </span>
            </a>
        </h3>
        <bar-chart :data="[['X-Small', 55], ['Small', 27], ['Hard', 22], ['Medium', 7]]"></bar-chart>

        <h3 class="boxed">Map of Players by Country</h3>
        <!-- https://github.com/paliari/v-autocomplete -->
        <geo-chart :library="{backgroundColor: '#EADBC4', datalessRegionColor: '#ded0ba'}" :data="[['United States', 44], ['Germany', 23], ['Brazil', 22]]"></geo-chart>

        <h3 class="boxed">Character Choice by Region</h3>
        <bar-chart :data="[['X-Small', 55], ['Small', 27], ['Hard', 22], ['Medium', 7]]"></bar-chart>

        <div class="split">
            <div>
                <h3 class="boxed">Servers by Platform</h3>
                <pie-chart :data="[['Blueberry', 44], ['Strawberry', 23]]"></pie-chart>
            </div>
            <div>
                <h3 class="boxed">Vanilla vs Modded Servers</h3>
                <pie-chart :data="[['Blueberry', 6], ['Strawberry', 15]]"></pie-chart>
            </div>
        </div>

        <div class="split">
            <div>
                <h3 class="boxed">Servers by Intent</h3>
                <bar-chart :data="[['X-Small', 55], ['Small', 27], ['Hard', 22], ['Medium', 7]]"></bar-chart>
            </div>
            <div>
                <h3 class="boxed">Servers by Season</h3>
                <bar-chart :data="[['X-Small', 55], ['Small', 27], ['Hard', 22], ['Medium', 7]]"></bar-chart>
            </div>
        </div>

    </main>
</template>

<script>
import axios from 'axios';

export default {
    methods: {
        get(endpoint) {
            return axios
                .get("http://localhost:3000" + endpoint);
        },
        transform(obj) {
            let arr = []
            for (var item in obj) {
                arr.push([item, obj[item]]);
            }
            return arr.sort(function(a, b) {
                return b[1] - a[1];
            });
        }
    },

    data() {
        return {
            characters: [],
        }
    },

    mounted() {
        this.get("/characters").then(resp => (this.characters = this.transform(resp.data)));
    }
}
</script>