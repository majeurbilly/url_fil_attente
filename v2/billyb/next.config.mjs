/** @type {import('next').NextConfig} */
const nextConfig = {
    async headers() {
        return [{ // correspondance de toutes les routes API
            source: "/api/list",
            headers: [
                {key: "Access-Control-Allow-Credentials", value: "true"}, {
                    key: "Access-Control-Allow-Origin", value: "https://10.100.2.130:3000/web/"
                }, {
                    key: "Access-Control-Allow-Methods", value: "GET,DELETE,POST,OPTIONS"
                }, {
                    key: "Access-Control-Allow-Headers", value: "X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, Origin, X-Requested-With, Content-Type, Accept, Authorization"
                },]
        }]
    }
};

export default nextConfig;
