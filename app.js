/**
 * Created by penguin on 2017/07/21.
 */
new Vue({
    el: '#app',
    data: {
        results: []
    },
    methods: {
        gets() {
            axios.get("/rules").then((response) => {
                console(response.data)
                this.results = response.data.results;
            }).catch((error) => { console.log(error); });
        }
    }
});