ELEMENT.locale(ELEMENT.lang.en)
new Vue({
    el: '#app',
    methods: {
        gets() {
            axios.get("/rules").then((response) => {
                console.log(response.data)
                Vue.set(this, "tableData", response.data)
            }).catch((error) => { console.log(error); });
        },
        get(id) {
            axios.get("/rule/"+id).then((response) => {
                console.log(response.data)
                Vue.set(this, "form", response.data)
            }).catch((error) => { console.log(error); });
        },
        post() {
            axios.post("/rule", this.$data.form).then((response) => {
                console.log(response.data)
                this.gets()
            }).catch((error) => { console.log(error); });
        },
        put(id) {
            axios.put("/rule/"+id, this.$data.form).then((response) => {
                console.log(response.data)
                this.gets()
            }).catch((error) => { console.log(error); });
        },
        change(id, index) {
            axios.put("/rule/"+id, {Enable: !this.$data.tableData[index].Enable}).then((response) => {
                console.log(response.data)
                this.gets()
            }).catch((error) => { console.log(error); });
        },
        handleDelete(index, row) {
            console.log(index, row);
            axios.delete("/rule/"+row.Id).then((response) => {
                console.log(response.data)
                this.gets()
            }).catch((error) => { console.log(error); });
        },
        dialogFormVisibleClean() {
            Vue.set(this, "form", {})
        },
        handleSelect(key, keyPath) {
            this.activeIndex = key
            console.log(key, keyPath);
        }
    },
    data: function() {
        return {
            activeIndex: '1',
            tableData: [],
            dialogFormVisible: false,
            dialogFormVisibleUpdate: false,
            targetId: '',
            form: {
                LocalAddr: '',
                SshAddr: '',
                RemoteAddr: ''
              },
            formLabelWidth: '120px'
        }
    },
    created: function() {
        this.gets()
    }
});