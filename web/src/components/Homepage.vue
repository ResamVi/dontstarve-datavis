<template>
    <main> 
        <h1>One two three</h1>
        <p>
            For a guide and recipes on how to configure / customize this project,<br>
            check out the
            <a href="https://cli.vuejs.org" target="_blank" rel="noopener">vue-cli documentation</a>.
        </p>
        <h2> Lorem ipsum whatever</h2> 
        <h3 class="boxed"> Yeah this is enough </h3>

        <line-chart :data="{'2017-01-01': 11, '2017-01-02': 6}"></line-chart>

    </main>
</template>