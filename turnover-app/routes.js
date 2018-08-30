//SPDX-License-Identifier: Apache-2.0

var box = require('./controller.js');

module.exports = function(app){

  app.get('/get_tuna/:id', function(req, res){
    box.get_tuna(req, res);
  });
  app.get('/add_box/:box', function(req, res){
    box.add_tuna(req, res);
  });
  app.get('/get_all_box', function(req, res){
    box.get_all_tuna(req, res);
  });
  app.get('/change_holder/:holder', function(req, res){
    box.change_holder(req, res);
  });
}
