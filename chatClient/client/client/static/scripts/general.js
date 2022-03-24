$(document).ready(function () {
    $('#nav').load('./navbar.html');
    $('#footer').load('./footer.html');
});

/**
 * 获取URL中的参数
 * @param variable URL中的参数名
 * @returns {string|boolean} 参数值
 */
function getQueryVariable(variable) {
    const query = window.location.search.substring(1);
    const vars = query.split("&");
    for (let i = 0; i < vars.length; i++) {
        let pair = vars[i].split("=");
        if (pair[0] === variable) {
            return pair[1];
        }
    }
    return false;
}

/**
 * 跳转到指定的企业页面
 * @param entId 企业ID号
 */
function redirectToEnt(entId) {
    window.location.href='../enterprise.html?entId=' + entId;
}
