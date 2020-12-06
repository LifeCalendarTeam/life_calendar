var dayButtons = new Vue({
    el: "#dayButtons",

    data: {
        days: []
    },

    methods: {
        colorToHex(color) {
            function compToHex(c) {
                let hex = c.toString(16);
                return hex.length === 1 ? "0" + hex : hex;
            }
            return "#" + compToHex(color[0]) + compToHex(color[1]) + compToHex(color[2]);
        },

        isDayFilled(day) {
            for(let filledDay of this.days) {
                if(Date.parse(filledDay.date) === Date.parse(day.toISOString())) {
                    return true;
                }
            }
            return false;
        },

        getFilledDayColor(day) {
            for(let filledDay of this.days) {
                if(Date.parse(filledDay.date) === Date.parse(day.toISOString())) {
                    return filledDay.average_color;
                }
            }
            return [169, 169, 169];
        },

        getDisplayedWeeks() {
            let closestSunday = new Date();
            closestSunday = new Date(Date.UTC(closestSunday.getFullYear(), closestSunday.getMonth(),
                closestSunday.getDate() + 7 - closestSunday.getDay()));
            let threeWeeks = [];
            for(let i = 2; i >= 0; --i) {
                let currentWeek = [];
                for(let j = 6; j >= 0; --j) {
                    currentWeek.push(new Date(closestSunday));
                    currentWeek[6 - j].setUTCDate(closestSunday.getUTCDate() - j - i * 7);
                }
                threeWeeks.push(currentWeek);
            }
            return threeWeeks;
        }
    },

    mounted() {
        axios.get("/api/days/brief").then((response) => {
            this.days = response.data;
        });
    }

});