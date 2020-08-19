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
                <h1>{{ serverCount }}</h1>
                <p>Servers analyzed</p>
            </div>
            <div>
                <h1>{{ playerCount }}</h1>
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
        <bar-chart :data="serverOrigin"></bar-chart>

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
        <bar-chart :data="playerOrigin"></bar-chart>

        <h3 class="boxed">Map of Players by Country</h3>
        <geo-chart :library="{backgroundColor: '#EADBC4', datalessRegionColor: '#ded0ba'}" :data="playerOrigin"></geo-chart>

        <h3 class="boxed">Character Choice by Region</h3>
        <div class="center">
            <autocomplete id="inputfield" :search="search" placeholder="United States" @submit="submit"></autocomplete>
        </div>

        <bar-chart :data="charactersOrigin"></bar-chart>

        <div class="split">
            <div>
                <h3 class="boxed">Servers by Platform</h3>
                <pie-chart :data="platformCount"></pie-chart>
            </div>
            <div>
                <h3 class="boxed">Vanilla vs Modded Servers</h3>
                <pie-chart :data="moddedCount"></pie-chart>
            </div>
        </div>

        <div class="split">
            <div>
                <h3 class="boxed">Servers by Intent</h3>
                <bar-chart :data="intentCount"></bar-chart>
            </div>
            <div>
                <h3 class="boxed">Servers by Season</h3>
                <bar-chart :data="seasonCount"></bar-chart>
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
                arr.push([this.capitalize(item), obj[item]]);
            }
            return arr.sort(function(a, b) {
                return b[1] - a[1];
            });
        },
        capitalize(string) {
            return string.charAt(0).toUpperCase() + string.slice(1);
        },

        search(input) {
            if (input.length < 1) { return [] }
            
            return this.countries.filter(country => {
                return country.toLowerCase().startsWith(input.toLowerCase())
            })
        },

        submit(input) {
            input = input.replace(" ", "%20");
            this.get("/characters/" + input).then(resp => (this.charactersOrigin = this.transform(resp.data)));
        }
    },

    data() {
        return {
            characters: [],
            countries: [],
            charactersOrigin: [],
            playerCount: 0,
            serverCount: 0,
            playerOrigin: [],
            serverOrigin: [],
            intentCount: [],
            platformCount: [],
            moddedCount: [],
            seasonCount: []
        }
    },

    mounted() {
        this.get("/countries").then(resp => (this.countries = resp.data));
        this.get("/characters").then(resp => (this.characters = this.transform(resp.data)));
        this.get("/attribute/intent").then(resp => (this.intentCount = this.transform(resp.data)));
        this.get("/attribute/platforms").then(resp => (this.platformCount = this.transform(resp.data)));
        this.get("/attribute/modded").then(resp => (this.moddedCount = this.transform(resp.data)));
        this.get("/attribute/season").then(resp => (this.seasonCount = this.transform(resp.data)));
        this.get("/characters/united%20states").then(resp => (this.charactersOrigin = this.transform(resp.data)));
        this.get("/count/players").then(resp => (this.playerCount = resp.data));
        this.get("/count/servers").then(resp => (this.serverCount = resp.data));
        this.get("/origin/players").then(resp => (this.playerOrigin = this.transform(resp.data)));
        this.get("/origin/servers").then(resp => (this.serverOrigin = this.transform(resp.data)));
    }
}
</script>

<style>
/* Override library CSS */
.autocomplete-input {
    padding: 9px !important;
    background-image: none;
    font-size: 10px;
    text-align: center;
}

.center {
    text-align: center;
}
</style>