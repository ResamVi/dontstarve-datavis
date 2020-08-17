<template>
    <main>
        <div id="text">
            <img src="../assets/logo-short.png" id="logo-short">
            <div>
                <p>
                    Welcome to my DST DataViz - a website recording and evaluating current and past character usage across different regions of players of Don't Starve Together.
                </p>
                <p>
                    This website is dedicated to players who have asked themselves "But who is the most played character now?". Along with answering this question more information
                    as far as Klei's servers of what is documented and allowed are shown like regional differences.
                </p>
                <p>
                    This project was created by <a href="https://resamvi.io/">ResamVi</a> and you can find the code <a href="https://github.com/ResamVi">here</a>.
                    The backend logic and this website are GPL licensed (feel free to modify, but keep it open for everyone). 
                    The project is still in its infancy with development started in August 2020. We may see in the future where this will be going.
                <p>
                    I'm always a bit swamped with university but otherwise would be happy about any contact and meeting everyone in-game. Don't hesitate to reach out:
                    <a href="https://steamcommunity.com/profiles/76561198969315128/">Steam</a>
                </p>
                <p>
                    Sincerely, ResamVi
                </p>
            </div>
        </div>
    </main>
</template>

<style>
main {
    margin-top: 10px;
}

#logo-short {
    float: left;
    width: 50%;
    height: 50%;
    margin: 10px;
}

@media only screen and (min-width: 1000px) {
  #logo-short {
    width: 380px !important;
    height: 238px !important;
  }
  main {
      margin: 0 25% !important;
  }
}
</style>