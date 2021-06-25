<template>
    <main> 
        <h1>Welcome!</h1>
        <p>
            DST DataViz records and visualizes current player preferences across regions
            on Klei's <a href="https://store.steampowered.com/app/322330/Dont_Starve_Together/">
            Don't Starve Together</a> game. <br/>
            Data generated {{ Math.round(age) }} minutes ago.<br/>
            Started tracking since: {{ started }}.
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

        <h3 class="boxed">Count of Characters being played</h3>
        <span style="float:right;">
            <input type="checkbox" v-model="includeModdedChars" @click="toggleModdedChar">
            <label for="modded-chars">Include Modded Characters</label>
        </span>
        <bar-chart :data="characters"></bar-chart>

        <h3 class="boxed">Amount of Servers by Country</h3>
        <bar-chart :data="countryCount"></bar-chart>

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
        <bar-chart :data="playerCountry"></bar-chart>

        <h3 class="boxed">Map of Players by Country</h3>
        <geo-chart :library="{backgroundColor: '#EADBC4', datalessRegionColor: '#ded0ba'}" :data="allPlayerCountry"></geo-chart>

        <h3 class="boxed">Character Choice by Country</h3>
        <div class="center">
            <autocomplete id="country-field" :search="searchCountry" placeholder="Enter Country here (Currently: China)" @submit="submitCountry"></autocomplete>
        </div>
        <bar-chart :data="charactersCountry"></bar-chart>

        <h3 class="boxed">Activity over time by Continent</h3>
        <line-chart :data="seriesContinents" />

        <h3 class="boxed">
            Count of Characters over time
            <a href="#" class="has-tooltip">[ ? ]
                <span class="tooltip tooltip-top">
                    Times are in Central European Standard Time (CEST) or UTC+2. If you are in PST subtract 9, in EDT subtract 6.
                </span>
            </a>
        </h3>
        <line-chart :data="seriesCharacter" />

        <h3 class="boxed">Highest Character Preferences by Country</h3>
        <div>
            <flag-row character="Wilson"         :data=wilson />
            <flag-row character="Willow"         :data=willow />
            <flag-row character="Wolfgang"       :data=wolfgang />
            <flag-row character="Wendy"          :data=wendy />
            <flag-row character="WX-78"          :data=wx />
            <flag-row character="Wickerbottom"   :data=wickerbottom />
            <flag-row character="Woodie"         :data=woodie />
            <flag-row character="Wes"            :data=wes />
            <flag-row character="Waxwell"        :data=waxwell />
            <flag-row character="Wigfrid"        :data=wigfrid />
            <flag-row character="Webber"         :data=webber />
            <flag-row character="Warly"          :data=warly />
            <flag-row character="Wormwood"       :data=wormwood />
            <flag-row character="Winona"         :data=winona />
            <flag-row character="Wortox"         :data=wortox />
            <flag-row character="Wurt"           :data=wurt />
            <flag-row character="Walter"         :data=walter />
        </div>

        <h3 class="boxed">Character Preference over time</h3>
        <autocomplete id="character-field" :search="searchCharacter" placeholder="Enter Character" @submit="submitCharacter"></autocomplete>
        <flag-column :character="characterInput" :data="seriesPreferences" />

        <h3 class="boxed">
            Player Preferences
            <a href="#" class="has-tooltip">[ ? ]
                <span class="tooltip tooltip-top">
                    Started tracking since: {{ started }}
                </span>
            </a>
        </h3>
        <div style="display: flex">
            <input type="text" id="player-field" name="player-field" class="autocomplete-input" placeholder="Enter your in-game name (case sensitive)" v-model="playerInput">
            <button @click="submitPlayer">Submit</button>
        </div>

        <div class="boxed">
            <div>
                <h3>Characters played (in hours)</h3>
                <column-chart :data="playerCharacters"></column-chart>
            </div>
        </div>

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
import FlagRow from '../components/FlagRow.vue';
import FlagColumn from '../components/FlagColumn.vue';

export default {
    components: {
        FlagRow,
        FlagColumn
    },
    methods: {
        get(endpoint) {
            return axios.get(process.env.VUE_APP_SERVER_ENDPOINT + endpoint);
        },
        
        searchCountry(input) {
            if (input.length < 1) { return [] }
            
            return this.countries.filter(country => {
                return country.toLowerCase().startsWith(input.toLowerCase())
            })
        },
        
        searchCharacter(input) {
            if (input.length < 1) { return [] }
            
            const characters = ["Wilson", "Willow", "Wolfgang", "Wendy", "WX-78", "Wickerbottom", "Woodie", "Wes", "Maxwell", "Wigfrid", "Webber", "Warly", "Wormwood", "Winona", "Wortox", "Wurt", "Walter"];

            return characters.filter(country => {
                return country.toLowerCase().startsWith(input.toLowerCase())
            })
        },

        submitCountry(input) {
            this.countryInput = input.replace(" ", "%20");
            this.get("/characters/" + this.countryInput).then(resp => (this.charactersCountry = resp.data));
        },

        submitCharacter(input) {
            this.characterInput = input;
            this.get("/series/preferences/" + input).then(resp => (this.seriesPreferences = resp.data));
        },

        submitPlayer() {
            this.get("/player/character/" + this.playerInput).then(resp => (this.playerCharacters = resp.data));
        },

        toggleModdedChar() {
            setTimeout(() => { 
                this.get("/characters?modded=" + this.includeModdedChars).then(resp => (this.characters = resp.data));
            }, 100);
        },

        
    },

    data() {
        return {
            age: 0,
            started: "",
            countryInput: "China",
            characterInput: "Wilson",
            playerInput: "",

            playerCharacters: [],

            playerCount: 0,
            serverCount: 0,
            characters: [],
            countries: [],
            charactersCountry: [],
            playerCountry: [],
            countryCount: [],
            intentCount: [],
            platformCount: [],
            moddedCount: [],
            seasonCount: [],
            allPlayerCountry: [],
            
            includeModdedChars: false,

            wilson: [],
            willow: [],
            wolfgang: [],
            wendy: [],
            wx: [],
            wickerbottom: [],
            woodie: [],
            wes: [],
            waxwell: [],
            wigfrid: [],
            webber: [],
            warly: [],
            wormwood: [],
            winona: [],
            wortox: [],
            wurt: [],
            walter: [],

            seriesCharacter: [],
            seriesPreferences: [],
            seriesContinents: [
                []
            ],
        }
    },

    mounted() {
        this.get("/meta/servers")               .then(resp => (this.serverCount = resp.data));
        this.get("/meta/players")               .then(resp => (this.playerCount = resp.data));
        this.get("/characters")                 .then(resp => (this.characters = resp.data));
        this.get("/count/country")              .then(resp => (this.countryCount = resp.data));
        this.get("/count/players")              .then(resp => (this.playerCountry = resp.data));
        this.get("/count/players?all=true")     .then(resp => (this.allPlayerCountry = resp.data));
        this.get("/meta/countries")             .then(resp => (this.countries = resp.data));
        this.get("/characters/China")           .then(resp => (this.charactersCountry = resp.data));
        this.get("/series/continents")          .then(resp => (this.seriesContinents = resp.data));
        this.get("/series/characters")          .then(resp => (this.seriesCharacter = resp.data));
        
        this.get("/characters/percentage/Wilson")       .then(resp => (this.wilson = resp.data));
        this.get("/characters/percentage/Willow")       .then(resp => (this.willow = resp.data));
        this.get("/characters/percentage/Wolfgang")     .then(resp => (this.wolfgang = resp.data));
        this.get("/characters/percentage/Wendy")        .then(resp => (this.wendy = resp.data));
        this.get("/characters/percentage/WX-78")        .then(resp => (this.wx = resp.data));
        this.get("/characters/percentage/Wickerbottom") .then(resp => (this.wickerbottom = resp.data));
        this.get("/characters/percentage/Woodie")       .then(resp => (this.woodie = resp.data));
        this.get("/characters/percentage/Wes")          .then(resp => (this.wes = resp.data));
        this.get("/characters/percentage/Maxwell")      .then(resp => (this.waxwell = resp.data));
        this.get("/characters/percentage/Wigfrid")      .then(resp => (this.wigfrid = resp.data));
        this.get("/characters/percentage/Webber")       .then(resp => (this.webber = resp.data));
        this.get("/characters/percentage/Warly")        .then(resp => (this.warly = resp.data));
        this.get("/characters/percentage/Wormwood")     .then(resp => (this.wormwood = resp.data));
        this.get("/characters/percentage/Winona")       .then(resp => (this.winona = resp.data));
        this.get("/characters/percentage/Wortox")       .then(resp => (this.wortox = resp.data));
        this.get("/characters/percentage/Wurt")         .then(resp => (this.wurt = resp.data));
        this.get("/characters/percentage/Walter")       .then(resp => (this.walter = resp.data));
        
        this.get("/count/platforms")            .then(resp => (this.platformCount = resp.data));
        this.get("/count/intent")               .then(resp => (this.intentCount = resp.data));
        this.get("/count/modded")               .then(resp => (this.moddedCount = resp.data));
        this.get("/count/season")               .then(resp => (this.seasonCount = resp.data));
        this.get("/meta/age")                   .then(resp => (this.age = resp.data));
        this.get("/meta/started")               .then(resp => (this.started = resp.data));
        
        this.get("/series/preferences/Wilson")          .then(resp => (this.seriesPreferences = resp.data));
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

.portrait {
    width: 6vw;
    height: 7vw;
    justify-self: center;
    align-self: center;
}
</style>
