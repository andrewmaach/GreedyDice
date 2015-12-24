module.exports = function(grunt) {

  // Project configuration.
  grunt.initConfig({
    shell: {
      build: {
          command: "go build ."
      },
      run: {
          command: "GreedyAI"
      }
    }
  });
  
  grunt.loadNpmTasks('grunt-shell');
  
  grunt.registerTask('build', 'Build GreedyAI', function() {
    var done = this.async();
    require('child_process').exec('go build .',
        function (error, stdout, stderr) {
            if (error !== null) {
             done(false);
            } else {
                done(true);
            }
        });
  });
  
  grunt.registerTask('build', ['shell:build']);
  grunt.registerTask('run', ['shell:run']);
  
  grunt.registerTask('default', ['build', "run"]);
  
  grunt.registerTask('test', ['build', "run"]);
  
}