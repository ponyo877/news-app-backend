db.article_col.aggregate([
    { $sort: { publishedAt: -1 } },
    { $group: {
        _id: { titles: "$titles", siteID: "$siteID" },
        total: { $sum: 1 },
        items: { $push: "$_id" }
    } },
    { $match: { total: { $gt: 1 } } },
]).toArray().then((docs) => {
    console.log(JSON.stringify(docs));
    var procs = [];
    for (var doc of docs) {
      doc.targets.shift();
      procs[procs.length] = db.article_col.deleteMany({
        _id: { $in: doc.targets }
      });
    }
    return Promise.all(procs);
  }).then((results) => {
    console.log("Remove dupulicate data.");
  }).catch((err) => {
    console.log(err);
  })
