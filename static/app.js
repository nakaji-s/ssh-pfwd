/**
 * Created by penguin on 2017/07/21.
 */
new Vue({
    el: '#app',
    data: {
    },
    methods: {
        gets() {
            axios.get("/rules").then((response) => {
                console.log(response.data)
                console.log(response.data[0].Priority)
            }).catch((error) => { console.log(error); });
        }
    }
});