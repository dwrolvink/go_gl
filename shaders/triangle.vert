#version 450
uniform float x;
uniform float y;
uniform float t;
in vec3 vp;
void main() {
    gl_Position = vec4(vp[0]+x+sin(0.01*t)/2.0, // x pos
                       vp[1]+y+cos(0.01*t)/2.0, // y pos
                       vp[2],                   // ?
                       cos(1.*(x+y))            // scale / z?
                  );
}