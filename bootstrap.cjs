// SPDX-License-Identifier: GPL-3.0
// ChunCord restart-safe updates. Maintained by https://github.com/chunma-san.
(() => {
    const fs = require("original-fs");
    const crypto = require("crypto");
    const target = __CHUNCORD_PATH__;
    const pending = target + ".pending";
    const receipt = pending + ".json";
    try {
        if (fs.existsSync(pending) && fs.existsSync(receipt)) {
            if (fs.lstatSync(receipt).size > 2048) throw new Error("Invalid update receipt");
            const info = JSON.parse(fs.readFileSync(receipt, "utf8"));
            const stat = fs.lstatSync(pending);
            if (info.version !== 1 || !stat.isFile() || stat.isSymbolicLink() || stat.size !== info.size || stat.size < 1 || stat.size > 128 * 1024 * 1024
                || !/^[a-f0-9]{64}$/.test(info.sha256) || !/^[a-f0-9]{40}$/.test(info.revision))
                throw new Error("Invalid staged update");
            const data = fs.readFileSync(pending);
            if (crypto.createHash("sha256").update(data).digest("hex") !== info.sha256)
                throw new Error("Staged update checksum mismatch");
            if (fs.existsSync(target)) fs.copyFileSync(target, target + ".previous");
            fs.renameSync(pending, target);
            fs.unlinkSync(receipt);
        }
    } catch (error) {
        // A bad candidate or a still-running process must leave Discord usable.
        try { fs.writeFileSync(target + ".update-error.log", String(error) + "\n"); } catch {}
    }
    require(target);
})();
