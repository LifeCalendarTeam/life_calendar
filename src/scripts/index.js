var day_buttons = new Vue({
    el:'#day_buttons',

    data: {
        days: []
    },

    methods: {
        colorToHex: function (color) {
            function compToHex(c) {
                let hex = c.toString(16);
                return hex.length === 1 ? "0" + hex : hex;
            }
            return "#" + compToHex(color[0]) + compToHex(color[1]) + compToHex(color[2]);
        },

        getDayUrl: function (day) {
            return "/" + day.getFullYear() + "_" + day.getMonth() + "_" + day.getDate();
        },

        isDayFilled: function (day) {
            for(let i = 0; i < this.days.length; ++i)
            {
                if(Date.parse(this.days[i].date) === Date.parse(day.toISOString()))
                    return true;
            }
            return false;
        },

        getFilledDayColor: function (day) {
            for(let i = 0; i < this.days.length; ++i)
            {
                if(Date.parse(this.days[i].date) === Date.parse(day.toISOString()))
                    return this.days[i].average_color;
            }
            return [169, 169, 169];
        },

        getDisplayedWeeks: function () {
            let closest_sunday = new Date();
            closest_sunday = new Date(Date.UTC(closest_sunday.getFullYear(), closest_sunday.getMonth(),
                closest_sunday.getDate() + 7 - closest_sunday.getDay()));
            let three_weeks = [];
            for(let i = 2; i >= 0; --i)
            {
                let current_week = [];
                for(let j = 6; j >= 0; --j)
                {
                    current_week.push(new Date(closest_sunday));
                    current_week[6 - j].setUTCDate(closest_sunday.getUTCDate() - j - i * 7);
                }
                three_weeks.push(current_week);
            }
            return three_weeks;
        }
    },
    mounted() {
        axios.get("/api/days/brief").then(response => {
            this.days = response.data;
        })
    }

})