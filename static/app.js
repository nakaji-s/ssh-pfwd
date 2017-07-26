new Vue({
    el: '#app',
    methods: {
        gets() {
            axios.get("/rules").then((response) => {
                console.log(response.data)
                Vue.set(this, "tableData", response.data)
            }).catch((error) => { console.log(error); });
        },
        put() {
            axios.put("/rule", this.$data.form).then((response) => {
                console.log(response.data)
                this.gets()
            }).catch((error) => { console.log(error); });
        },
        handleEdit(index, row) {
            console.log(index, row);
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
        }
    },
    data: function() {
        return {
            tableData: [],
            dialogFormVisible: false,
            form: {
                localaddr: '',
                sshaddr: '',
                remoteaddr: ''
              },
            formLabelWidth: '120px'
        }
    },
    created: function() {
        this.gets()
    }
});