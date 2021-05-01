#version 450
uniform float x;
uniform float y;
uniform float z;
uniform float scale;
uniform float t;
in vec3 vp;
void main() {
    gl_Position = vec4(vp[0]+x, // x pos
                       vp[1]+y, // y pos
                       vp[2],                   // ?
                       1./scale            
                  );
}