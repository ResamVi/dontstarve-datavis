<template>
    <main>
        <div id="text">
            <img src="../assets/logo-short.png" id="logo-short">
            <div>
                <p>
                    This website is dedicated to players who have asked themselves "But who is the most played character now?" Along with answering this question more information
                    (as far as what's documented and allowed) are shown - this covers mostly players, servers and server metadata, regional differences, platforms.
                </p>
                <p>
                    This project was created by <a href="https://resamvi.io/">ResamVi</a> and you can find the code <a href="https://github.com/ResamVi">here</a>.
                    The backend logic and this website are GPL licensed (feel free to modify, but keep it open for everyone). 
                    The project is a side-gig I did for fun in August 2020 and is currently in maintenance mode.
                </p>
                <p>
                    You can write me a <a href="mailto:julien.midedji@gmail.com">mail</a> if you need something 
                </p>
                <p>
                    Sincerely, ResamVi
                </p>
            </div>
        </div>
        <h3 class="boxed">Previous discussion on:</h3>
        <div class="split">
            <div>
                <a href="https://www.reddit.com/r/dontstarve/comments/idxxel/infographic_on_the_most_used_character_in_total/">
                    <img src="../assets/reddit-logo.png" width="120px"><br/>
                    Reddit <br/>
                    /r/dontstarve Thread
                </a>
            </div>
            <div>
                <a href="https://forums.kleientertainment.com/forums/topic/121116-most-used-character-in-total-by-region-and-more-infographic/">
                    <img src="../assets/klei-logo.png" width="120px"><br/>
                    Klei Forums<br/>
                    General Discussion Thread
                </a>
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